package rip

import "fmt"

type RouteMethodError struct {
	route *Route
	old   string
	new   string
}

func (rem *RouteMethodError) Error() string {
	return fmt.Sprintf("overwrote method %v with %v in %v", rem.old, rem.new, rem.route.Template())
}

type RouteMissingMethodError struct {
	route *Route
}

func (rem *RouteMissingMethodError) Error() string {
	return fmt.Sprintf("route has no method defined in %v", rem.route.Template())
}
