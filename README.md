# kit-blog
a microservice blog with go-kit



kit new service users
    change and add meths in service

kit g s users  --dmw -t grpc
    generate grpcs and default middleware for service and endpoint



## WATCHER
watcher -run github.com/emadghaffari/kit-blog/notificator/cmd  -watch github.com/emadghaffari/kit-bg/notificator

watcher -run github.com/emadghaffari/kit-blog/users/cmd  -watch github.com/emadghaffari/kit-bg/users

watcher -run github.com/emadghaffari/kit-blog/comments/cmd  -watch github.com/emadghaffari/kit-bg/comments

## configs:
for config management we use vault and consul

REPO: https://github.com/testdrivenio/vault-consul-docker