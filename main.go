
package main

import (
    "context"
    "fmt"

    "github.com/dev404ai/steampipe-plugin-yandexcloud/yandexcloud"
    "github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

var version string

func main() {
    plugin.Serve(&plugin.ServeOpts{
        PluginFunc: func(ctx context.Context) *plugin.Plugin {
            p := yandexcloud.Plugin()
            if version != "" {
                p.Name = fmt.Sprintf("%s@%s", p.Name, version)
            }
            return p
        },
    })
}
