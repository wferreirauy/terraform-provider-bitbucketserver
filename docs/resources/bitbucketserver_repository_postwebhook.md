# Resource: bitbucketserver_repository_postwebhook

Manage a repository level Post Webhook. Extends what Bitbucket does every time a repository or pull request occurs, for example when code is pushed or a pull request is merged.
\
Bitbucket Post Webhooks allow sending JSON data to any HTTP or HTTPS address when any of the selected events occur.
\
Marketplace Plugin Page - [Post Webhooks for Bitbucket](https://marketplace.atlassian.com/apps/1215474/post-webhooks-for-bitbucket?tab=overview&hosting=datacenter).
\
API Reference - [Atlassian Bitbucket Post Webhook API](https://help.moveworkforward.com/BPW/atlassian-bitbucket-post-webhook-api#AtlassianBitbucketPostWebhookAPI-Getconfigurations).

## Example Usage

```hcl
resource "bitbucketserver_project" "myproj" {
  key  = "MYPROJ"
  name = "my-project"
}

resource "bitbucketserver_repository" "repo" {
  project = bitbucketserver_project.myproj.key
  name    = "repo"
}

resource "bitbucketserver_repository_postwebhook" "jenkins" {
  project             = bitbucketserver_project.myproj.key
  repository          = bitbucketserver_repository.repo.slug
  title               = "Jenkins"
  webhook_url         = "https://jenkins.example.com/bitbucket-hook"
  commiters_to_ignore = "john.doe,jane.doe"
  branches_to_ignore  = "release/.*"
  enabled             = true
  repo_push           = true
  pr_merged           = true
}
```

## Argument Reference

* `project` - Required. Project Key the repository is contained within.
* `repository` - Required. Slug of the repository to which the post webhook will belong.
* `title` - Required. Title of the post webhook.
* `webhook_url` - Required. The URL of the post webhook.
* `commiters_to_ignore` - Optional. Comma separated list of usernames. Commits from these users do not trigger this hook.
* `branches_to_ignore` - Optional. Regex for branches. Commits on these branches do not trigger this hook.
* `active` - Optional. Enable or disable the webhook. Default: true
* `repo_push` - Optional. Event On push. Default: false
* `branch_deleted` - Optional. Event Branch deleted. Default: false
* `branch_created` - Optional. Event Branch created. Default: false
* `tag_created` - Optional. Event Tag created. Default: false
* `build_status` - Optional. Event Build status. Default: false
* `repo_mirror_synced` - Optional. Event Repository mirror synchronized. Default: false
* `pr_declined` - Optional. Event Pull request declined. Default: false
* `pr_rescoped` - Optional. Event Pull request re-scoped. Default: false
* `pr_merged` - Optional. Event Pull request merged. Default: false
* `pr_reopened` - Optional. Event Pull request re-opened. Default: false
* `pr_updated` - Optional. Event Pull request description updated. Default: false
* `pr_created` - Optional. Event Pull request created. Default: false
* `pr_commented` - Optional. Event Pull request commented. Default: false

## Attribute Reference

* `webhook_id` - The post webhook id.

## Import

Import a user reference using the project key, repository name and post webhook title.

```
terraform import bitbucketserver_repository_postwebhook.jenkins MYPROJ/repo/Jenkins
```
