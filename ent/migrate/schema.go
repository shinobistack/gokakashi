// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AgentTasksColumns holds the columns for the "agent_tasks" table.
	AgentTasksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "status", Type: field.TypeString, Default: "pending"},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "agent_id", Type: field.TypeInt},
		{Name: "scan_id", Type: field.TypeUUID},
	}
	// AgentTasksTable holds the schema information for the "agent_tasks" table.
	AgentTasksTable = &schema.Table{
		Name:       "agent_tasks",
		Columns:    AgentTasksColumns,
		PrimaryKey: []*schema.Column{AgentTasksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "agent_tasks_agents_agent_tasks",
				Columns:    []*schema.Column{AgentTasksColumns[3]},
				RefColumns: []*schema.Column{AgentsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "agent_tasks_scans_agent_tasks",
				Columns:    []*schema.Column{AgentTasksColumns[4]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// AgentsColumns holds the columns for the "agents" table.
	AgentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "status", Type: field.TypeString, Default: "connected"},
		{Name: "workspace", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "server", Type: field.TypeString, Nullable: true},
		{Name: "last_seen", Type: field.TypeTime},
	}
	// AgentsTable holds the schema information for the "agents" table.
	AgentsTable = &schema.Table{
		Name:       "agents",
		Columns:    AgentsColumns,
		PrimaryKey: []*schema.Column{AgentsColumns[0]},
	}
	// IntegrationTypesColumns holds the columns for the "integration_types" table.
	IntegrationTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "display_name", Type: field.TypeString, Unique: true},
	}
	// IntegrationTypesTable holds the schema information for the "integration_types" table.
	IntegrationTypesTable = &schema.Table{
		Name:       "integration_types",
		Columns:    IntegrationTypesColumns,
		PrimaryKey: []*schema.Column{IntegrationTypesColumns[0]},
	}
	// IntegrationsColumns holds the columns for the "integrations" table.
	IntegrationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "type", Type: field.TypeString},
		{Name: "config", Type: field.TypeJSON},
		{Name: "integration_type_integrations", Type: field.TypeString, Nullable: true},
	}
	// IntegrationsTable holds the schema information for the "integrations" table.
	IntegrationsTable = &schema.Table{
		Name:       "integrations",
		Columns:    IntegrationsColumns,
		PrimaryKey: []*schema.Column{IntegrationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "integrations_integration_types_integrations",
				Columns:    []*schema.Column{IntegrationsColumns[4]},
				RefColumns: []*schema.Column{IntegrationTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PoliciesColumns holds the columns for the "policies" table.
	PoliciesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "image", Type: field.TypeJSON},
		{Name: "scanner", Type: field.TypeString},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
		{Name: "trigger", Type: field.TypeJSON, Nullable: true},
		{Name: "notify", Type: field.TypeJSON, Nullable: true},
	}
	// PoliciesTable holds the schema information for the "policies" table.
	PoliciesTable = &schema.Table{
		Name:       "policies",
		Columns:    PoliciesColumns,
		PrimaryKey: []*schema.Column{PoliciesColumns[0]},
	}
	// PolicyLabelsColumns holds the columns for the "policy_labels" table.
	PolicyLabelsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "key", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
		{Name: "policy_id", Type: field.TypeUUID},
	}
	// PolicyLabelsTable holds the schema information for the "policy_labels" table.
	PolicyLabelsTable = &schema.Table{
		Name:       "policy_labels",
		Columns:    PolicyLabelsColumns,
		PrimaryKey: []*schema.Column{PolicyLabelsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "policy_labels_policies_policy_labels",
				Columns:    []*schema.Column{PolicyLabelsColumns[3]},
				RefColumns: []*schema.Column{PoliciesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// ScanLabelsColumns holds the columns for the "scan_labels" table.
	ScanLabelsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "key", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
		{Name: "scan_id", Type: field.TypeUUID},
	}
	// ScanLabelsTable holds the schema information for the "scan_labels" table.
	ScanLabelsTable = &schema.Table{
		Name:       "scan_labels",
		Columns:    ScanLabelsColumns,
		PrimaryKey: []*schema.Column{ScanLabelsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "scan_labels_scans_scan_labels",
				Columns:    []*schema.Column{ScanLabelsColumns[3]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// ScanNotifiesColumns holds the columns for the "scan_notifies" table.
	ScanNotifiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "hash", Type: field.TypeString},
		{Name: "scan_id", Type: field.TypeUUID},
	}
	// ScanNotifiesTable holds the schema information for the "scan_notifies" table.
	ScanNotifiesTable = &schema.Table{
		Name:       "scan_notifies",
		Columns:    ScanNotifiesColumns,
		PrimaryKey: []*schema.Column{ScanNotifiesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "scan_notifies_scans_scan_notifications",
				Columns:    []*schema.Column{ScanNotifiesColumns[2]},
				RefColumns: []*schema.Column{ScansColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// ScansColumns holds the columns for the "scans" table.
	ScansColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "status", Type: field.TypeString, Default: "scan_pending"},
		{Name: "image", Type: field.TypeString},
		{Name: "scanner", Type: field.TypeString},
		{Name: "notify", Type: field.TypeJSON, Nullable: true},
		{Name: "report", Type: field.TypeJSON, Nullable: true},
		{Name: "integration_id", Type: field.TypeUUID},
		{Name: "policy_id", Type: field.TypeUUID},
	}
	// ScansTable holds the schema information for the "scans" table.
	ScansTable = &schema.Table{
		Name:       "scans",
		Columns:    ScansColumns,
		PrimaryKey: []*schema.Column{ScansColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "scans_integrations_scans",
				Columns:    []*schema.Column{ScansColumns[6]},
				RefColumns: []*schema.Column{IntegrationsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "scans_policies_scans",
				Columns:    []*schema.Column{ScansColumns[7]},
				RefColumns: []*schema.Column{PoliciesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AgentTasksTable,
		AgentsTable,
		IntegrationTypesTable,
		IntegrationsTable,
		PoliciesTable,
		PolicyLabelsTable,
		ScanLabelsTable,
		ScanNotifiesTable,
		ScansTable,
	}
)

func init() {
	AgentTasksTable.ForeignKeys[0].RefTable = AgentsTable
	AgentTasksTable.ForeignKeys[1].RefTable = ScansTable
	IntegrationsTable.ForeignKeys[0].RefTable = IntegrationTypesTable
	PolicyLabelsTable.ForeignKeys[0].RefTable = PoliciesTable
	ScanLabelsTable.ForeignKeys[0].RefTable = ScansTable
	ScanNotifiesTable.ForeignKeys[0].RefTable = ScansTable
	ScansTable.ForeignKeys[0].RefTable = IntegrationsTable
	ScansTable.ForeignKeys[1].RefTable = PoliciesTable
}
