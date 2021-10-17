package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/brunoa19/ketch-terraform-provider/client/v1beta1"
	"github.com/google/go-containerregistry/pkg/name"
	registryv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type App struct {
	// +immutable
	Name      string `json:"name"`
	Image     string `json:"image"`
	Framework string `json:"framework"`
	// +optional
	Cname []string `json:"cnames,omitempty"`
	// +optional
	Ports []int `json:"ports"`
	// +optional
	Units int64 `json:"units"`
	// +optional
	Processes []*ProcessParameters `json:"processes,omitempty"`
	// +optional
	RoutingSettings *RoutingSettings `json:"routing_settings,omitempty"`
	// +optional
	Version int64 `json:"version"`
}

// ProcessParameters defines process parameters
type ProcessParameters struct {
	Cmd  []string `json:"cmd"`
	Name string   `json:"name"`
}

// RoutingSettings defines routing settings
type RoutingSettings struct {
	Weight int64 `json:"weight"`
}

func NewApp(input *v1beta1.App) *App {
	var deployment v1beta1.AppDeploymentSpec
	var ports []int
	var processes []*ProcessParameters
	if len(input.Spec.Deployments) > 0 {
		deployment = input.Spec.Deployments[0]
		for _, port := range deployment.ExposedPorts {
			ports = append(ports, port.Port)
		}

		for _, p := range deployment.Processes {
			processes = append(processes, &ProcessParameters{
				Name: p.Name,
				Cmd:  p.Cmd,
			})
		}
	}

	return &App{
		Name:      input.ObjectMeta.Name,
		Image:     deployment.Image,
		Framework: input.Spec.Framework,
		Cname:     input.Spec.Ingress.Cnames,
		Ports:     ports,
		Units:     int64(input.Spec.DeploymentsCount),
		Processes: processes,
		RoutingSettings: &RoutingSettings{
			int64(deployment.RoutingSettings.Weight),
		},
		Version: int64(deployment.Version),
	}
}

func newExposedPort(port string) (*v1beta1.ExposedPort, error) {
	parts := strings.SplitN(port, "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid port: " + port)
	}
	portInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New("invalid port: " + port)
	}
	return &v1beta1.ExposedPort{
		Port:     portInt,
		Protocol: strings.ToUpper(parts[1]),
	}, nil
}

//nolint:gocyclo
func (a *App) convertToKetchApp() *v1beta1.App {
	cfg, err := getImageConfig(a.Image)
	if err != nil {
		log.Println("#### GetImageConfig:ERR ", err)
	}

	var cmd []string
	if cfg != nil {
		cmd = make([]string, 0, len(cfg.Config.Entrypoint))
		cmd = append(cmd, cfg.Config.Entrypoint...)
		cmd = append(cmd, cfg.Config.Cmd...)
	}

	app := &v1beta1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name: a.Name,
		},
		Spec: v1beta1.AppSpec{
			Framework: a.Framework,
			Deployments: []v1beta1.AppDeploymentSpec{
				{
					Image: a.Image,
				},
			},
		},
	}

	if len(a.Ports) > 0 {
		for _, port := range a.Ports {
			app.Spec.Deployments[0].ExposedPorts = append(app.Spec.Deployments[0].ExposedPorts, v1beta1.ExposedPort{
				Port:     port,
				Protocol: "TCP",
			})
		}
	} else if cfg != nil {
		var exposedPorts []v1beta1.ExposedPort
		for port := range cfg.Config.ExposedPorts {
			exposedPort, err := newExposedPort(port)
			if err == nil {
				exposedPorts = append(exposedPorts, *exposedPort)
			}
		}

		app.Spec.Deployments[0].ExposedPorts = exposedPorts
	}

	if len(app.Spec.Deployments[0].ExposedPorts) == 0 {
		app.Spec.Deployments[0].ExposedPorts = append(app.Spec.Deployments[0].ExposedPorts, v1beta1.ExposedPort{
			Port:     8000,
			Protocol: "TCP",
		})
	}

	app.Spec.Ingress.GenerateDefaultCname = true
	if len(a.Cname) > 0 {
		app.Spec.Ingress.Cnames = append(app.Spec.Ingress.Cnames, a.Cname...)
	}

	// default 1
	if a.Units == 0 {
		app.Spec.DeploymentsCount = 1
	} else {
		app.Spec.DeploymentsCount = int(a.Units)
	}

	if len(a.Processes) > 0 {
		for _, p := range a.Processes {
			app.Spec.Deployments[0].Processes = append(app.Spec.Deployments[0].Processes, v1beta1.ProcessSpec{
				Name: p.Name,
				Cmd:  p.Cmd,
			})
		}
	} else if len(cmd) > 0 {
		app.Spec.Deployments[0].Processes = append(app.Spec.Deployments[0].Processes, v1beta1.ProcessSpec{
			Name: "web",
			Cmd:  cmd,
		})
	}

	if a.RoutingSettings != nil && a.RoutingSettings.Weight > 0 {
		app.Spec.Deployments[0].RoutingSettings.Weight = uint8(a.RoutingSettings.Weight)
	} else {
		app.Spec.Deployments[0].RoutingSettings.Weight = 100
	}

	if a.Version > 0 {
		app.Spec.Deployments[0].Version = v1beta1.DeploymentVersion(a.Version)
	}

	return app
}

// Wrapf wraps error and supplies the line and the file where the error occurred.
func Wrapf(err error, fmtStr string, params ...interface{}) error {
	_, fl, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(fmtStr, params...)
	return fmt.Errorf("message: %q; error: \"%w\"; file: %s; line: %d", msg, err, path.Base(fl), line)
}

func getImageConfig(imageName string) (*registryv1.ConfigFile, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, Wrapf(err, "failed to parse reference for image %q", imageName)
	}
	var options []remote.Option
	img, err := remote.Image(ref, options...)
	if err != nil {
		return nil, Wrapf(err, "could not get config for image %q", imageName)
	}
	return img.ConfigFile()
}

func (c *Client) GetApp(ctx context.Context, name string) (*App, error) {
	app, err := c.getApp(ctx, name)
	if err != nil {
		return nil, err
	}

	return NewApp(app), nil
}

func (c *Client) DeleteApp(ctx context.Context, name string) error {
	app, err := c.getApp(ctx, name)
	if err != nil {
		return err
	}

	err = c.kube.Delete(ctx, app)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getApp(ctx context.Context, name string) (*v1beta1.App, error) {
	app := &v1beta1.App{}
	err := c.kube.Get(ctx, types.NamespacedName{Name: name}, app)
	if err != nil {
		log.Println("ERR:", err.Error())
		return nil, err
	}

	return app, nil
}

func (c *Client) CreateApp(ctx context.Context, input *App) error {
	app := input.convertToKetchApp()
	return c.kube.Create(ctx, app)
}

func (c *Client) UpdateApp(ctx context.Context, input *App) error {
	app, err := c.getApp(ctx, input.Name)
	if err != nil {
		return err
	}

	updates := input.convertToKetchApp()
	app.Spec = updates.Spec
	return c.kube.Update(ctx, app)
}
