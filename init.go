package eureka

import core "github.com/procyon-projects/procyon-core"

func init() {
	// properties
	core.Register(newClientProperties)
	core.Register(newInstanceProperties)
	// instance info provider
	core.Register(newDefaultInstanceInfoProvider)
}
