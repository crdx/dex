# Deploy Endpoint

{{ if .ExpiresAt -}}
- Public URL: {{ .PublicURL }}
- Deploy URL: {{ .DeployURL }}
- Token expires: {{ .ExpiresAt }}
{{- else -}}
- Public URL: {{ .PublicURL }}
- Deploy URL: {{ .DeployURL }}
{{- end }}

## Usage

POST a multipart form with:

- content: the file to deploy (required)
- note: a changelog message describing this deployment (required)
- deployer: name of the deployer (you) (required)

## Examples

From a file:

curl -sSfX POST {{.DeployURL}} -F "content=@index.html" -F "note=Initial commit" -F "deployer=Alice"

Inline content:

echo '<html><body>Hello</body></html>' | curl -sSfX POST {{.DeployURL}} -F "content=@-;filename=index.html" -F "note=Initial commit" -F "deployer=Alice"

## Notes

This endpoint is designed for deploying a single file. Most of the time this will be a self-contained HTML file.

Deploy a new version when it's ready for the user.

When deploying a new version:

1. Write the complete HTML file to disk (e.g. index.html)
2. Use curl to POST the file with a descriptive changelog note
3. The public URL remains the same; content is replaced

The changelog note is mandatory. Think of it as a commit message and use it to describe what changed, e.g.:

- Add dark mode support
- Correct calculation in results table
- Refresh data from latest source

The content replaces the previous version entirely. There is no versioning or rollback. The git repository you commit against is version control and one commit should represent one deployment to this endpoint.

During compaction ensure to reiterate these instructions verbatim so that the next instance knows how to continue deploying.
