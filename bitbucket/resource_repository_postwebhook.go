package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type PostWebhook struct {
	ID                 int    `json:"id,omitempty"`
	Title              string `json:"title,omitempty"`
	URL                string `json:"url,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	CommittersToIgnore string `json:"committersToIgnore,omitempty"`
	BranchesToIgnore   string `json:"branchesToIgnore,omitempty"`
	TagCreated         bool   `json:"tagCreated,omitempty"`
	BranchDeleted      bool   `json:"branchDeleted,omitempty"`
	BranchCreated      bool   `json:"branchCreated,omitempty"`
	RepoPush           bool   `json:"repoPush,omitempty"`
	PrDeclined         bool   `json:"prDeclined,omitempty"`
	PrRescoped         bool   `json:"prRescoped,omitempty"`
	PrMerged           bool   `json:"prMerged,omitempty"`
	PrReopened         bool   `json:"prReopened,omitempty"`
	PrUpdated          bool   `json:"prUpdated,omitempty"`
	PrCreated          bool   `json:"prCreated,omitempty"`
	PrCommented        bool   `json:"prCommented,omitempty"`
	BuildStatus        bool   `json:"buildStatus,omitempty"`
	RepoMirrorSynced   bool   `json:"repoMirrorSynced,omitempty"`
}

type PostWebhookListResponse struct {
	Values []PostWebhook `json:"$"`
}

func resourceRepositoryPostWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryPostWebhookCreate,
		Update: resourceRepositoryPostWebhookUpdate,
		Read:   resourceRepositoryPostWebhookRead,
		Delete: resourceRepositoryPostWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"webhook_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"webhook_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"committers_to_ignore": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"branches_to_ignore": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_created": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"branch_deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"branch_created": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"repo_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_declined": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_rescoped": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_merged": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_reopened": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_updated": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_created": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pr_commented": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"build_status": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"repo_mirror_synced": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceRepositoryPostWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	id := d.Get("webhook_id").(int)
	webhook := newPostWebhookFromResource(d)

	request, err := json.Marshal(webhook)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("/rest/webhook/1.0/projects/%s/repos/%s/configurations/%d",
		project,
		repository,
		id,
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	return resourceRepositoryPostWebhookRead(d, m)
}

func resourceRepositoryPostWebhookCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient

	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	webhook := newPostWebhookFromResource(d)

	request, err := json.Marshal(webhook)

	if err != nil {
		return err
	}

	res, err := client.Post(fmt.Sprintf("/rest/webhook/1.0/projects/%s/repos/%s/configurations",
		project,
		repository,
	), bytes.NewBuffer(request))

	if err != nil {
		return err
	}

	var webhookResponse PostWebhook

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &webhookResponse)

	if err != nil {
		return err
	}

	_ = d.Set("webhook_id", webhookResponse.ID)

	d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("project").(string), d.Get("repository").(string), d.Get("title").(string)))
	return resourceRepositoryPostWebhookRead(d, m)
}

func resourceRepositoryPostWebhookRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id != "" {
		parts := strings.Split(id, "/")
		if len(parts) == 3 {
			_ = d.Set("project", parts[0])
			_ = d.Set("repository", parts[1])
			_ = d.Set("title", parts[2])
		} else {
			return fmt.Errorf("incorrect ID format, should match `project/repository/title`")
		}
	}

	var err error = getRepositoryPostWebhookFromList(d, m)

	if err != nil {
		return err
	}

	return nil
}

func resourceRepositoryPostWebhookDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BitbucketServerProvider).BitbucketClient
	_, err := client.Delete(fmt.Sprintf("/rest/webhook/1.0/projects/%s/repos/%s/configurations/%d",
		d.Get("project").(string),
		d.Get("repository").(string),
		d.Get("webhook_id").(int)))

	return err
}

func newPostWebhookFromResource(d *schema.ResourceData) (Hook *PostWebhook) {
	webhook := &PostWebhook{
		Title:              d.Get("title").(string),
		URL:                d.Get("webhook_url").(string),
		Enabled:            d.Get("enabled").(bool),
		CommittersToIgnore: d.Get("committers_to_ignore").(string),
		BranchesToIgnore:   d.Get("branches_to_ignore").(string),
		TagCreated:         d.Get("tag_created").(bool),
		BranchDeleted:      d.Get("branch_deleted").(bool),
		BranchCreated:      d.Get("branch_created").(bool),
		RepoPush:           d.Get("repo_push").(bool),
		PrDeclined:         d.Get("pr_declined").(bool),
		PrRescoped:         d.Get("pr_rescoped").(bool),
		PrMerged:           d.Get("pr_merged").(bool),
		PrReopened:         d.Get("pr_reopened").(bool),
		PrUpdated:          d.Get("pr_updated").(bool),
		PrCreated:          d.Get("pr_created").(bool),
		PrCommented:        d.Get("pr_commented").(bool),
		BuildStatus:        d.Get("build_status").(bool),
		RepoMirrorSynced:   d.Get("repo_mirror_synced").(bool),
	}

	return webhook
}

func getRepositoryPostWebhookFromList(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	repository := d.Get("repository").(string)
	title := d.Get("title").(string)

	client := m.(*BitbucketServerProvider).BitbucketClient

	resp, err := client.Get(fmt.Sprintf("/rest/webhook/1.0/projects/%s/repos/%s/configurations",
		project,
		repository,
	))

	if err != nil {
		return err
	}

	var webhookListResponse PostWebhookListResponse

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&webhookListResponse)

	if err != nil {
		return err
	}

	for _, webhook := range webhookListResponse.Values {
		if webhook.Title == title {
			_ = d.Set("webhook_id", webhook.ID)
			_ = d.Set("webhook_url", webhook.URL)
			_ = d.Set("enabled", webhook.Enabled)
			_ = d.Set("committers_to_ignore", webhook.CommittersToIgnore)
			_ = d.Set("branches_to_ignore", webhook.BranchesToIgnore)
			_ = d.Set("tag_created", webhook.TagCreated)
			_ = d.Set("branch_deleted", webhook.BranchDeleted)
			_ = d.Set("branch_created", webhook.BranchCreated)
			_ = d.Set("repo_push", webhook.RepoPush)
			_ = d.Set("pr_declined", webhook.PrDeclined)
			_ = d.Set("pr_rescoped", webhook.PrRescoped)
			_ = d.Set("pr_merged", webhook.PrMerged)
			_ = d.Set("pr_reopened", webhook.PrReopened)
			_ = d.Set("pr_updated", webhook.PrUpdated)
			_ = d.Set("pr_created", webhook.PrCreated)
			_ = d.Set("pr_commented", webhook.PrCommented)
			_ = d.Set("build_status", webhook.BuildStatus)
			_ = d.Set("repo_mirror_synced", webhook.RepoMirrorSynced)
			return nil
		}
	}

	return fmt.Errorf("incorrect ID format, should match `project/repository/title`")
}
