variable "env_vars" { type = object({
  new_york_times_key      = string
  openapi_key                 = string
  twitter_api_key             = string
  twitter_api_key_secret      = string
  twitter_access_token        = string
  twitter_access_token_secret = string
}) }

variable "project_name" {
  type = string
  default = "content-generator"
}