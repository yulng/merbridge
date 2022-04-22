/*
Copyright © 2022 Merbridge Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cniplugin

import (
	"testing"

	"istio.io/istio/cni/pkg/plugin"
)

func TestIgnorePod(t *testing.T) {
	type args struct {
		namespace string
		name      string
		pod       *plugin.PodInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "uninjected pod",
			args: args{
				pod: &plugin.PodInfo{
					Containers: []string{"foo", "bar"},
				},
			},
			want: true,
		},
		{
			name: "injected pod without sidecar status",
			args: args{
				pod: &plugin.PodInfo{
					Containers: []string{"foo", "istio-proxy"},
				},
			},
			want: true,
		},
		{
			name: "injected pod",
			args: args{
				pod: &plugin.PodInfo{
					Containers: []string{"foo", "istio-proxy"},
					Annotations: map[string]string{
						sidecarStatusKey: "whatever",
					},
				},
			},
			want: false,
		},
		{
			name: "injected pod with envoy disabled",
			args: args{
				pod: &plugin.PodInfo{
					Containers: []string{"foo", "istio-proxy"},
					Annotations: map[string]string{
						sidecarStatusKey: "whatever",
					},
					ProxyEnvironments: map[string]string{
						"DISABLE_ENVOY": "true",
					},
				},
			},
			want: true,
		},
		{
			name: "injected pod with envoy enabled",
			args: args{
				pod: &plugin.PodInfo{
					Containers: []string{"foo", "istio-proxy"},
					Annotations: map[string]string{
						sidecarStatusKey: "whatever",
					},
					ProxyEnvironments: map[string]string{
						"DISABLE_ENVOY": "false",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ignorePod(tt.args.namespace, tt.args.name, tt.args.pod); got != tt.want {
				t.Errorf("ignorePod() = %v, want %v", got, tt.want)
			}
		})
	}
}
