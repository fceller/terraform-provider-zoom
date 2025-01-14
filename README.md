This terraform provider allows to perform Create ,Read ,Update, Delete, Import and Deactivate Zoom Users.

## Requirements

* [Go](https://golang.org/doc/install) 1.16 <br>
* [Terraform](https://www.terraform.io/downloads.html) 0.13.x <br/>
* [Zoom](https://zoom.us/) Pro/Premium account
* [Zoom API Documentations](https://marketplace.zoom.us/docs/api-reference/zoom-api/users)

## Application Account

***This provider can only be successfully tested on a premium paid zoom account.*** <br><br>

### Setup

1. Create a zoom account with paid subscription (PRO Plan/Business Account). (https://zoom.us/)<br>

### API Authentication

1. Go to [Zoom Marketplace](https://marketplace.zoom.us/)<br>
2. Click on `Build App`. For our purpose we need to make a JWT App. <br>
3. Create a server-to-server application.
4. This app will provide us with the ZOOM_ACCOUNT_ID, ZOOM_CLIENT_ID, and ZOOM_CLIENT_SECRET which will be needed
   to configure our provider and make request. <br>

## Building the Provider

1. Clone the repository, add all the dependencies and create a vendor directory that contains all dependencies. For
   this, run the following commands:

```cd terraform-provider-zoom
go mod init terraform-provider-zoom
go mod tidy
go mod vendor
```

## Managing terraform plugins

1. Run the following command to create a vendor subdirectory which will comprise of all provider dependencies. <br>

```
%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
``` 

Command:

```bash
mkdir -p %APPDATA%/terraform.d/plugins/hashicorp.com/edu/zoom/0.2.0/[OS_ARCH]
```

For eg. `mkdir -p %APPDATA%/terraform.d/plugins/hashicorp.com/edu/zoom/0.2.0/windows_amd64`<br>

2. Run `go build -o terraform-provider-zoom.exe`. This will save the binary (`.exe`) file in the main/root
   directory. <br>
3. Run this command to move this binary file to appropriate location.

 ```
 move terraform-provider-zoom.exe %APPDATA%\terraform.d\plugins\hashicorp.com\edu\zoom\0.2.0\[OS_ARCH]
 ``` 

[OR]

1. Manually move the file from current directory to destination directory <br>

## Working with terraform

### Application Credential Integration in terraform

1. Add `terraform` block and `provider` block as shown in [example usage](#example-usage).
2. Get a pair of credentials: zoom_api_secret and zoom_api_secret. For this, visit https://marketplace.zoom.us/.
3. Assign the above credentials to the respective field in the `provider` block.

### Basic Terraform Commands

1. `terraform init` - To initialize a working directory containing Terraform configuration files.
2. `terraform plan` - To create an execution plan. Displays the changes to be done.
3. `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the changes.

### Create User

1. Add the user email, first name, last name, status, license_type, deartment, job_title, location in the respective
   field as shown in [example usage](#example-usage).
2. Run the basic terraform commands.<br>
3. On successful execution, sends an account setup mail to user.<br>

### Update the user

1. Update the data of the user in the `resource` block as show in [example usage](#example-usage) and run the basic
   terraform commands to update user. User is not allowed to update `email`.

2. Update the `status` of User from `active` to `inactive` or viceversa and run the basic terraform commands.

### Read the User Data

Add `data` and `output` blocks as shown in the [example usage](#example-usage) and run the basic terraform commands.

### Delete the user

Delete the `resource` block of the user and run `terraform apply`.

#### Import a User Data

1. Write manually a `resource` configuration block for the user as shown in [example usage](#example-usage). Imported
   user will be mapped to this block.
2. Run the command `terraform import zoom_user.user1 [EMAIL_ID]` to import user.
3. Run `terraform plan`, if output shows `0 to add, 0 to change and 0 to destroy` user import is successful, otherwise
   recheck the user data in `resource` block with user data in Zoom website.

## Example Usage <a id="example-usage"></a>

```terraform
terraform {
  required_providers {
    zoom = {
      version >= "1.0.3"
      source  = "CleverTap/zoom"
    }
  }
}

provider "zoom" {
  zoom_account_id = "[ZOOM_ACCOUNT_ID]"
  zoom_client_id = "[ZOOM_CLIENT_ID]"
  zoom_client_secret = "[ZOOM_CLIENT_SECRET]"
  timeout_minutes = 3
}

resource "zoom_user" "user1" {
  email      = "useremail@domain.com"
  first_name = "Dummyfirst"
  last_name  = "Dummylast"
  status = "active"
  license_type = 1
  department = "DevOps"
  job_title = "Engineer"
  location   =  "India"
}

data "zoom_user" "user1" {
  email = "user@domain.com"
}

output "user1" {
  value = data.zoom_user.user1
}
```

## Argument Reference

* `zoom_account_id`(Required, string)     - The Zoom Account ID. This may also be set via the `"ZOOM_ACCOUNT_ID"` environment
  variable.
* `zoom_client_id`(Required, string)  - The Zoom Client ID. This may also be set via the `"ZOOM_CLIENT_ID"`
  environment variable.
* `zoom_client_secret`(Required, string)  - The Zoom Client secret. This may also be set via the `"ZOOM_CLIENT_SECRET"`
  environment variable.
* `timeout_minutes` (Optional, int)    - The duration for which retries to be performed when an API request fails with
  API Rate limit error. This may also be set via the `"ZOOM_TIMEOUT_MINUTES"` environment variable. Default value is 2.
* `email`(Required, string)            - The email id associated with the user account.
* `first_name`(Required, string)       - First name of the User.
* `last_name`(Required, string)        - Last Name / Family Name / Surname of the User.
* `status`(Optional, string)           - User account activation status i.e., `active`, `inactive`.
* `license_type`(Required, integer)    - User account type i.e.,(1=Basic, 2=License, 3=On-prem)
* `company` (Optional, string)         - User's company
* `job_title`(Optional, string)        - Job title of the particular user.
* `department`(Optional, string)       - Department of the particular user.
* `location`(Optional, string)         - Department of the particular user.
* `pmi`(Optional, integer)             - Personal Meeting ID of the user.
* `use_pmi`(Optional, bool)            - If true means use PMI for instant meetings.
* `timezone`(Optional, string)         - Time zone ID of user. For values
  refer [here](https://marketplace.zoom.us/docs/api-reference/other-references/abbreviation-lists#timezones).
* `vanity_name`(Optional, string)      - Personal Meeting Room Name.
* `host_key`(Optional, string)         - Host Key of user.
* `group_id`(Optional, string)         - Unique identifier of group to which user to be added.
* `manager`(Optional, string)          - Name or Email of Manager of the User.
* `phone_numbers`(Optional, set of maps) - Maps contain keys code, country, number, label.
* `language`(Optional, string)         - Language of User.
* `cms_user_id`(Optional, string)      - CMS User ID of User.
* `pronouns`(Optional, string)         - User's pronouns.
* `pronouns_option`(Optional, int)     - User's display pronouns setting.
* `role_name`(Computed, string)        - Current role of the user i.e., (Admin,Member).
