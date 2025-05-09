package kubernetes.admission.pod

# Deny privileged containers
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    c := input.object.spec.containers[_]
    c.securityContext.privileged
    msg := sprintf("Privileged container '%s' is not allowed", [c.name])
}

# Deny containers with hostPID, hostIPC, or hostNetwork
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    input.object.spec.hostPID
    msg := "Pod hostPID is not allowed"
}

deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    input.object.spec.hostIPC
    msg := "Pod hostIPC is not allowed"
}

deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    input.object.spec.hostNetwork
    msg := "Pod hostNetwork is not allowed"
}

# Deny containers that mount sensitive host paths
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    volume := input.object.spec.volumes[_]
    hostpath := volume.hostPath
    not allowed_hostpath(hostpath.path)
    msg := sprintf("HostPath volume '%s' mountPath '%s' is not allowed", [volume.name, hostpath.path])
}

# Define allowed host paths (empty means none are allowed)
allowed_hostpath(path) = false {
    # By default, no host paths are allowed
    # You can define exceptions here if needed
    false
}

# Require containers have resource limits
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    c := input.object.spec.containers[_]
    not c.resources.limits
    msg := sprintf("Container '%s' does not have resource limits", [c.name])
}

# Warn about missing liveness probes
warn[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    c := input.object.spec.containers[_]
    not c.livenessProbe
    msg := sprintf("Container '%s' does not have a liveness probe", [c.name])
}

# Warn about missing readiness probes
warn[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    c := input.object.spec.containers[_]
    not c.readinessProbe
    msg := sprintf("Container '%s' does not have a readiness probe", [c.name])
}

# Disallow latest tags in container images
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    c := input.object.spec.containers[_]
    endswith(c.image, ":latest")
    msg := sprintf("Container '%s' uses an image with the 'latest' tag (%s)", [c.name, c.image])
}

# Require specific labels
deny[msg] {
    input.operation == "CREATE"
    input.kind == "Pod"
    required_labels := ["app", "environment", "owner"]
    provided_labels := {label | input.object.metadata.labels[label]}
    missing := required_labels - provided_labels
    count(missing) > 0
    msg := sprintf("Pod is missing required labels: %v", [missing])
} 