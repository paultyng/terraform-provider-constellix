package constellix

import (
	"github.com/Constellix/constellix-go-client/models"
)

// DomainAttributes extends the models.DomainAttributes.
type DomainAttributes struct {
	models.DomainAttributes
	Disabled bool `json:"disabled"`
}

// DomainAttributesV4 contains the domain attributes aligned with API v4.
type DomainAttributesV4 struct {
	Enabled bool `json:"enabled"`
}
