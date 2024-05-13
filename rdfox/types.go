package rdfox

import "context"

type (
	CreateConnection          func(ctx context.Context, server, protocol, role, password, datastore string) (string, error)
	DeleteConnection          func(ctx context.Context, server, protocol, role, password, datastore, connectionID string) error
	CreateCursor              func(ctx context.Context, server, protocol, role, password, datastore, connectionID, query string) (string, error)
	DeleteCursor              func(ctx context.Context, server, protocol, role, password, datastore, connectionID, cursorID string) error
	ReadWithCursor            func(ctx context.Context, server, protocol, role, password, datastore, connectionID, cursorID string, limit int, gotData func(Triples)) error
	ImportAxioms              func(ctx context.Context, protocol, server, role, password, datastore, srcGraph, dstGraph string) error
	GetRoles                  func(ctx context.Context, server, protocol, role, password string) ([]string, error)
	CreateRoles               func(ctx context.Context, server, protocol, role, password, newRoleName, newRolePassword string) error
	DeleteRole                func(ctx context.Context, server, protocol, role, password, roleToDelete string) error
	GrantDatastorePrivileges  func(ctx context.Context, server, protocol, role, password, targetRole, datastore, resource, accessTypes string) error
	RevokeDatastorePrivileges func(ctx context.Context, server, protocol, role, password, targetRole, datastore, resource, accessTypes string) error
	ListPrivileges            func(ctx context.Context, server, protocol, role, password, targetRole string) (map[string][]string, error)
	GetStats                  func(ctx context.Context, server, protocol, role, password, datastore string) (Statistics, error)
)

type (
	Statistics map[string]map[string]interface{}
	Triples    map[string]map[string][]string
)
