package config

import (
	"github.com/spf13/cobra"
)

//const (
//	ServiceLabel string = "app.kubernetes.io/name=tekton-results-api"
//)

type Options struct {
	Config   Config
	NoPrompt bool
}

func SetCommand() *cobra.Command {
	o := &Options{}
	return &cobra.Command{
		Use:   "set",
		Short: "Set results config",
		//RunE: func(cmd *cobra.Command, args []string) error {
		//	cfg, err := askConfig()
		//	if err != nil {
		//		return err
		//	}
		//	fmt.Println("Token:", cfg.APIToken)
		//	return nil
		//},
		RunE: o.Run,
	}
}

func (o *Options) Run(_ *cobra.Command, args []string) (err error) {
	o.Config, err = NewConfig()
	if err != nil {
		return
	}
	return o.Config.Set(true)
}

//func askConfig() (*cli.Extension, error) {
//	var cfg cli.Extension
//
//	host, err := getHost()
//	if err != nil {
//		return nil, err
//	}
//	// Ask for API URL
//	hostPrompt := &survey.Select{
//		Message: "Enter the API URL:",
//		Options: host,
//	}
//	err = survey.AskOne(hostPrompt, &cfg.Host)
//	if err != nil {
//		return nil, err
//	}
//
//	// Ask for API Token
//	authToken, err := getAuthToken()
//	if err != nil {
//		return nil, err
//	}
//	tokenPrompt := &survey.Input{
//		Message: "Enter the API Token:",
//		Default: authToken,
//	}
//	err = survey.AskOne(tokenPrompt, &cfg.APIToken)
//	if err != nil {
//		return nil, err
//	}
//
//	return &cfg, nil
//}

//func setExtension() error {
//	kubeconfigPath := clientcmd.RecommendedHomeFile
//	// Load kubeConfig
//	config, err := clientcmd.LoadFromFile(kubeconfigPath)
//	if err != nil {
//		return fmt.Errorf("failed to load kubeconfig: %w", err)
//	}
//
//	// Get the current context
//	contextName := config.CurrentContext
//	if contextName == "" {
//		return fmt.Errorf("no current context set in kubeconfig")
//	}
//
//	ctx, ok := config.Contexts[contextName]
//	if !ok {
//		return fmt.Errorf("context %q not found in kubeconfig", contextName)
//	}
//
//	// Ensure extensions map exists
//	if ctx.Extensions == nil {
//		ctx.Extensions = make(map[string]runtime.Object)
//
//	// Update or create the Tekton Results extension
//	ctx.Extensions["tekton-results"] = api.NewExtension(map[string]interface{}{
//		"api-path": apiPath,
//		"host":     host,
//		"token":    token,
//	})
//
//	// Persist changes using ModifyConfig
//	err = clientcmd.ModifyConfig(configAccess, *config, true)
//	if err != nil {
//		return fmt.Errorf("failed to update kubeconfig: %w", err)
//	}
//
//
//}
//
//func getAuthToken() (string, error) {
//	kubeconfigPath := clientcmd.RecommendedHomeFile
//
//	// Load kubeConfig
//	config, err := clientcmd.LoadFromFile(kubeconfigPath)
//	if err != nil {
//		return "", fmt.Errorf("failed to load kubeconfig: %w", err)
//	}
//
//	// Get the current context
//	contextName := config.CurrentContext
//	if contextName == "" {
//		return "", fmt.Errorf("no current context set in kubeconfig")
//	}
//
//	// Get context details
//	context, exists := config.Contexts[contextName]
//	if !exists {
//		return "", fmt.Errorf("context %q not found in kubeconfig", contextName)
//	}
//
//	// Get user credentials
//	authInfo, exists := config.AuthInfos[context.AuthInfo]
//	if !exists {
//		return "", fmt.Errorf("user %q not found in kubeconfig", context.AuthInfo)
//	}
//
//	// Check if token is present
//	if authInfo.Token != "" {
//		return authInfo.Token, nil
//	}
//
//	return "", fmt.Errorf("no authentication token found in kubeconfig")
//}
//
//// getK8sClients initializes Kubernetes core and OpenShift Route clients
//func getK8sClients() (*corev1.CoreV1Client, *routev1.RouteV1Client, error) {
//	var config *rest.Config
//	var err error
//
//	kubeconfigPath := clientcmd.RecommendedHomeFile
//	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
//		config, err = rest.InClusterConfig()
//	} else {
//		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
//	}
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to load kubeconfig: %w", err)
//	}
//
//	coreClient, err := corev1.NewForConfig(config)
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create Kubernetes core client: %w", err)
//	}
//
//	routeClient, err := routev1.NewForConfig(config)
//	if err != nil {
//		return nil, nil, fmt.Errorf("failed to create OpenShift Route client: %w", err)
//	}
//
//	return coreClient, routeClient, nil
//}
//
//// getHost retrieves the OpenShift Route URL based on the labeled service
//func getHost() ([]string, error) {
//	coreClient, routeClient, err := getK8sClients()
//	if err != nil {
//		return nil, err
//	}
//
//	// Find the service with the label selector
//	services, err := coreClient.Services("").List(context.Background(), metav1.ListOptions{
//		LabelSelector: ServiceLabel,
//	})
//	if err != nil || len(services.Items) == 0 {
//		return nil, fmt.Errorf("no service found with label %q", ServiceLabel)
//	}
//	var routes []string
//	for _, svc := range services.Items {
//		// Find the route that maps to this service
//		routeList, err := routeClient.Routes(svc.Namespace).List(context.Background(), metav1.ListOptions{})
//		if err != nil {
//			return nil, fmt.Errorf("failed to list routeList: %w", err)
//		}
//		for _, route := range routeList.Items {
//			if route.Spec.To.Name == svc.Name {
//				port := route.Spec.Port.TargetPort
//				for _, p := range svc.Spec.Ports {
//					if p.Port == port.IntVal || p.Name == port.StrVal {
//						routes = append(routes, route.Spec.Host)
//					}
//				}
//			}
//		}
//	}
//	return routes, nil
//}
