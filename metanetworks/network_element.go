package metanetworks

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	networkElementsEndpoint string = "/v1/network_elements"
)

type NetworkElement struct {
	Aliases       []string       `json:"aliases,omitempty" meta_api:"read_only"`
	CreatedAt     string         `json:"created_at,omitempty" meta_api:"read_only"`
	Description   string         `json:"description"`
	DNSName       string         `json:"dns_name,omitempty" meta_api:"read_only"`
	Enabled       bool           `json:"enabled,omitempty" meta_api:"read_only"`
	ExpiresAt     string         `json:"expires_at,omitempty" meta_api:"read_only"`
	ID            string         `json:"id,omitempty" meta_api:"read_only"`
	MappedService string         `json:"mapped_service,omitempty"`
	MappedSubnets []string       `json:"mapped_subnets,omitempty"`
	ModifiedAt    string         `json:"modified_at,omitempty" meta_api:"read_only"`
	Name          string         `json:"name"`
	NetID         int64          `json:"net_id,omitempty" meta_api:"read_only"`
	OrgID         string         `json:"org_id,omitempty" meta_api:"read_only"`
	OwnerID       string         `json:"owner_id,omitempty" meta_api:"read_only"`
	Platform      string         `json:"platform,omitempty" meta_api:"read_only"`
	Type          string         `json:"type,omitempty" meta_api:"read_only"`
	Version       int64          `json:"version,omitempty" meta_api:"read_only"`
	MappedDomains []MappedDomain `json:"mapped_domains,omitempty"`
}

func (c *Client) GetNetworkElement(elementID string) (*NetworkElement, error) {
	var networkElement NetworkElement
	err := c.Read(networkElementsEndpoint+"/"+elementID+"?expand=true", &networkElement)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Network Element from Get: %s", networkElement.MappedDomains)
	return &networkElement, nil
}

func (c *Client) UpdateNetworkElement(networkElementID string, networkElement *NetworkElement) (*NetworkElement, error) {
	resp, err := c.Update(networkElementsEndpoint+"/"+networkElementID, *networkElement)
	if err != nil {
		return nil, err
	}
	updatedNetworkElement, _ := resp.(*NetworkElement)

	log.Printf("Returning Network Element from Update: %s", updatedNetworkElement.ID)
	return updatedNetworkElement, nil
}

func (c *Client) CreateNetworkElement(networkElement *NetworkElement) (*NetworkElement, error) {
	resp, err := c.Create(networkElementsEndpoint, *networkElement)
	if err != nil {
		return nil, err
	}

	createdNetworkElement, ok := resp.(*NetworkElement)
	if !ok {
		return nil, errors.New("Object returned from API was not a Network Element Pointer")
	}

	log.Printf("Returning Network Element from Create: %s", createdNetworkElement.ID)
	return createdNetworkElement, nil
}

func (c *Client) DeleteNetworkElement(networkElementID string) error {
	err := c.Delete(networkElementsEndpoint + "/" + networkElementID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetNetworkElementAlias(networkElementID string, alias string) (*NetworkElement, error) {
	resp, err := c.Request(networkElementsEndpoint+"/"+networkElementID+"/aliases/"+alias, "PUT", nil, "")
	if err != nil {
		return nil, err
	}
	var networkElement NetworkElement
	err = json.Unmarshal(resp, &networkElement)
	if err != nil {
		return nil, err
	}

	return &networkElement, nil
}

func (c *Client) DeleteNetworkElementAlias(networkElementID string, alias string) (*NetworkElement, error) {
	resp, err := c.Request(networkElementsEndpoint+"/"+networkElementID+"/aliases/"+alias, "DELETE", nil, "")
	if err != nil {
		return nil, err
	}
	var networkElement NetworkElement
	err = json.Unmarshal(resp, &networkElement)
	if err != nil {
		return nil, err
	}
	return &networkElement, nil
}

func (c *Client) SetNetworkElementMappedDomains(networkElementID string, name string, mappedDomain *MappedDomain) (*MappedDomain, error) {
	resp, err := c.SetMappedDomain(networkElementsEndpoint+"/"+networkElementID+"/mapped_domains/"+name, *mappedDomain)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Network Element Mapped Domains from Create: %s", resp.Name)
	return resp, nil
}

func (c *Client) DeleteNetworkElementMappedDomains(networkElementID string, name string) error {
	err := c.Delete(networkElementsEndpoint + "/" + networkElementID + "/mapped_domains/" + name)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetNetworkElementTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		tagsMapInterface := d.Get("tags").(map[string]interface{})
		tagsMapString := make(map[string]string)
		for key, value := range tagsMapInterface {
			tagsMapString[key] = value.(string)
		}

		id := d.Id()
		err := c.UpdateTags(networkElementsEndpoint+"/"+id+"/tags", tagsMapString)
		if err != nil {
			return err
		}
	}

	return nil
}
