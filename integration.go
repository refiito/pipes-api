package main

import (
	"strings"
)

type (
	Integration struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		Link       string  `json:"link"`
		Image      string  `json:"image"`
		Pipes      []*Pipe `json:"pipes"`
		AuthURL    string  `json:"auth_url,omitempty"`
		AuthType   string  `json:"auth_type,omitempty"`
		Authorized bool    `json:"authorized"`
	}
)

var availableIntegration = map[string][]string{
	"basecamp": {"users", "projects", "todolists", "todos"},
}

var availableAuthorizations = map[string]string{
	"basecamp": "oauth",
}

var availableImages = map[string]string{
	"basecamp": "/images/logo-basecamp.png",
}

var availableDescriptions = map[string]string{
	"basecamp:users":     "Basecamp users will be imported as Toggl users. Existing users are matched by e-mail.",
	"basecamp:projects":  "Basecamp projects will be imported as Toggl projects. Existing projects are matched by name.",
	"basecamp:todolists": "Basecamp todolists will be imported as Toggl tasks. Existing tasks are matched by name.",
	"basecamp:todos":     "Basecamp todos will be imported as Toggl tasks. Existing tasks are matched by name.",
}

var automaticOptions = map[string]bool{
	"basecamp:users":     false,
	"basecamp:projects":  true,
	"basecamp:todolists": true,
	"basecamp:todos":     true,
}

var premiumOptions = map[string]bool{
	"basecamp:users":     false,
	"basecamp:projects":  false,
	"basecamp:todolists": true,
	"basecamp:todos":     true,
}

func NewIntegration(serviceName string) *Integration {
	integration := Integration{
		ID:       serviceName,
		Link:     "http://support.toggl.com/basecamp",
		Name:     strings.Title(serviceName),
		Image:    availableImages[serviceName],
		AuthType: availableAuthorizations[serviceName],
		AuthURL:  getAuthURL(serviceName),
	}
	return &integration
}

func workspaceIntegrations(workspaceID int) ([]*Integration, error) {
	// FIXME: if authorizations, workspace pipes, pipe statues
	// don't block each others loading, load all 3 at the same time.

	authorizations, err := loadAuthorizations(workspaceID)
	if err != nil {
		return nil, err
	}

	workspacePipes, err := loadPipes(workspaceID)
	if err != nil {
		return nil, err
	}

	pipeStatuses, err := loadPipeStatuses(workspaceID)
	if err != nil {
		return nil, err
	}

	var integrations []*Integration
	for serviceID, pipeIDs := range availableIntegration {
		integration := NewIntegration(serviceID)
		integration.Authorized = authorizations[serviceID]
		for _, pipeID := range pipeIDs {
			key := pipesKey(serviceID, pipeID)
			pipe := workspacePipes[key]
			if pipe == nil {
				pipe = NewPipe(workspaceID, serviceID, pipeID)
			}
			pipe.PipeStatus = pipeStatuses[key]
			pipe.Premium = premiumOptions[key]
			pipe.Description = availableDescriptions[key]
			pipe.AutomaticOption = automaticOptions[key]
			integration.Pipes = append(integration.Pipes, pipe)
		}
		integrations = append(integrations, integration)
	}
	return integrations, nil
}
