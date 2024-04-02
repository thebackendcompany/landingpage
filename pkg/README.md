# TheBackendCompany Server

Backend application serving lead generation from application forms.

##### Footnotes:

- The frontend code in cmd/server/static, is generated using nextjs
- The template expects to have a {{.crsf_token}} field in the contact form.
- There can be multiple integrations for syncing data, right now its google sheets.


## Integrations
The store to save email and other leads is pluggable. The immediate candidates are:
- Google Sheets
- Redis
- RDBMS (Postgres)


## Prequisites

- Sessions are used to manage **csrf-tokens** to prevent random tom-dick-harry from submitting form
- No html user inputs are allowed in this API, Incase you add one, make sure to **santize** inputs to prevent **XSS**
- Some of these store integration requires some upfront work. As listed below.


#### Creds storage

There is no storage layer for the credentials file right now. Its stored with the code, encrypted with
a **master key**. The **master key** is not shared. Without the master key, that file is unrecoverable.

**Encrypt**

```shell
cat ~/dir/creds.json | go run -encrypt config/gcloud.env.enc -key $MASTER_KEY
```
__this is a necessary step__


**Decrypt**

```shell
go run -decrypt config/gcloud.env.enc -key $MASTER_KEY
```

**Please make sure to use the proper file name in the `config/*.env` files**



#### Google Sheets Integrations

(The links have not be hyperlinked, Copy the URL -> Right Click -> Open Link)

**Step 1** Enabling service

- https://console.cloud.google.com/projectcreate , create a New Project.
- https://console.cloud.google.com/apis/library/sheets.googleapis.com, enable the google sheets api
- Once you enable, Create Security Credentials and chose Auth Data (one that doesn't require use interaction).
    - The UI keeps changing, but basically you need to create Service Account.
    - There is no need to Grant any IAM policies to this service account. Leave the optional fields.

**Step 2** Getting Keys

- Once Service Account is created, You will see a listing, with an email like `<service-account-name-or-id>@projectname.iam.gserviceaccount.com`.
    - Copy the email
    - Get the Keys in JSON format (from Manage Keys in dropown)
    - If key is not present, it will ask you to create one. Go ahead an create it.


**Step 3** Setting up Spreadsheet

- Go to https://sheets.google.com , and create your spreadsheet.
- Share the spreadsheet with the Service Account email __you just copied from above__.
- For the first row, write the column names.
- Copy the google's `sheet id` from the URL
    - A google spreadsheet url looks like `https://docs.google.com/spreadsheets/d/1iC7hTkwBlz_LL7bFP7kZFlpGBzDFGKeGlofOISw_oUs/edit#gid=0`,
