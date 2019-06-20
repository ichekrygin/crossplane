package v1alpha1

import (
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// CloudSQL Instance API to Type mapping
// Each mapped Property comes with set of two helper functions:
// - newProperty(from *sqladmin.Property) *Property - creates a new mapped property from the api equivalent
// - saveProperty(from *Property) *sqladmin.Propert - creates a new api property from the type equivalent

// AclEntry: An entry for an Access Control list.
type AclEntry struct {
	// ExpirationTime: The time when this access control entry expires in
	// RFC 3339 format, for example 2012-11-15T16:19:00.094Z.
	ExpirationTime string `json:"expirationTime,omitempty"`

	// Name: An optional label to identify this entry.
	Name string `json:"name,omitempty"`

	// Value: The whitelisted value for the access control list.
	Value string `json:"value,omitempty"`
}

func newAclEntry(from *sqladmin.AclEntry) *AclEntry {
	if from == nil {
		return nil
	}

	return &AclEntry{
		ExpirationTime: from.ExpirationTime,
		Name:           from.Name,
		Value:          from.Value,
	}
}

func saveAclEntry(from *AclEntry) *sqladmin.AclEntry {
	if from == nil {
		return nil
	}
	return &sqladmin.AclEntry{
		ExpirationTime: from.ExpirationTime,
		Name:           from.Name,
		Value:          from.Value,
	}
}

// BackupConfiguration: Database instance backup configuration.
type BackupConfiguration struct {
	// BinaryLogEnabled: Whether binary log is enabled. If backup
	// configuration is disabled, binary log must be disabled as well.
	BinaryLogEnabled bool `json:"binaryLogEnabled,omitempty"`

	// Enabled: Whether this configuration is enabled.
	Enabled bool `json:"enabled,omitempty"`

	// ReplicationLogArchivingEnabled: Reserved for future use.
	ReplicationLogArchivingEnabled bool `json:"replicationLogArchivingEnabled,omitempty"`

	// StartTime: Start time for the daily backup configuration in UTC
	// timezone in the 24 hour format - HH:MM.
	StartTime string `json:"startTime,omitempty"`
}

func newBackupConfiguration(from *sqladmin.BackupConfiguration) *BackupConfiguration {
	if from == nil {
		return nil
	}
	return &BackupConfiguration{
		BinaryLogEnabled:               from.BinaryLogEnabled,
		Enabled:                        from.Enabled,
		ReplicationLogArchivingEnabled: from.ReplicationLogArchivingEnabled,
		StartTime:                      from.StartTime,
	}
}

func saveBackupConfiguration(from *BackupConfiguration) *sqladmin.BackupConfiguration {
	if from == nil {
		return nil
	}
	return &sqladmin.BackupConfiguration{
		BinaryLogEnabled:               from.BinaryLogEnabled,
		Enabled:                        from.Enabled,
		ReplicationLogArchivingEnabled: from.ReplicationLogArchivingEnabled,
		StartTime:                      from.StartTime,
	}
}

// DatabaseFlags: Database flags for Cloud SQL instances.
type DatabaseFlags struct {
	// Name: The name of the flag. These flags are passed at instance
	// startup, so include both server options and system variables for
	// MySQL. Flags should be specified with underscores, not hyphens. For
	// more information, see Configuring Database Flags in the Cloud SQL
	// documentation.
	Name string `json:"name,omitempty"`

	// Value: The value of the flag. Booleans should be set to on for true
	// and off for false. This field must be omitted if the flag doesn't
	// take a value.
	Value string `json:"value,omitempty"`
}

func newDatabaseFlags(from *sqladmin.DatabaseFlags) *DatabaseFlags {
	return &DatabaseFlags{
		Name:  from.Name,
		Value: from.Value,
	}
}

func saveDatabaseFlags(from *DatabaseFlags) *sqladmin.DatabaseFlags {
	return &sqladmin.DatabaseFlags{
		Name:  from.Name,
		Value: from.Value,
	}
}

// DatabaseInstanceFailoverReplica: The name and status of the failover
// replica. This property is applicable only to Second Generation
// instances.
type DatabaseInstanceFailoverReplica struct {
	// Available: The availability status of the failover replica. A false
	// status indicates that the failover replica is out of sync. The master
	// can only failover to the falover replica when the status is true.
	Available bool `json:"available,omitempty"`

	// Name: The name of the failover replica. If specified at instance
	// creation, a failover replica is created for the instance. The name
	// doesn't include the project ID. This property is applicable only to
	// Second Generation instances.
	Name string `json:"name,omitempty"`
}

func newDatabaseInstanceFailoverReplica(from *sqladmin.DatabaseInstanceFailoverReplica) *DatabaseInstanceFailoverReplica {
	if from == nil {
		return nil
	}
	return &DatabaseInstanceFailoverReplica{
		Available: from.Available,
		Name:      from.Name,
	}
}

func saveDatabaseInstanceFailoverReplica(from *DatabaseInstanceFailoverReplica) *sqladmin.DatabaseInstanceFailoverReplica {
	if from == nil {
		return nil
	}
	return &sqladmin.DatabaseInstanceFailoverReplica{
		Available: from.Available,
		Name:      from.Name,
	}
}

// DatabaseInstanceSpec: A Cloud SQL instance resource input.
type DatabaseInstanceSpec struct {
	// DatabaseVersion: The database engine type and version.
	// MySQL Second Generation instances: MYSQL_5_7 (default) or MYSQL_5_6.
	// PostgreSQL instances: POSTGRES_9_6 MySQL First Generation instances:
	// MYSQL_5_6 (default) or MYSQL_5_5
	//
	// The databaseVersion field can not be changed after instance creation.
	DatabaseVersion string `json:"databaseVersion"`

	// InstanceType: The instance type. This can be one of the following.
	// CLOUD_SQL_INSTANCE: A Cloud SQL instance that is not replicating from a master.
	// READ_REPLICA_INSTANCE: A Cloud SQL instance configured as a read-replica.
	InstanceType string `json:"instanceType,omitempty"`

	// MasterInstanceName: The name of the instance which will act as master
	// in the replication setup.
	MasterInstanceName string `json:"masterInstanceName,omitempty"`

	// Region: The geographical region. Can be us-central (FIRST_GEN
	// instances only), us-central1 (SECOND_GEN instances only), asia-east1
	// or europe-west1. Defaults to us-central or us-central1 depending on
	// the instance type (First Generation or Second Generation).
	//
	// The region can not be changed after instance creation.
	Region string `json:"region"`

	// ReplicaConfiguration: Configuration specific to failover replicas and read replicas.
	ReplicaConfiguration *ReplicaConfiguration `json:"replicaConfiguration,omitempty"`

	// Settings: The user Settings.
	Settings *Settings `json:"settings,omitempty"`
}

func newDatabaseInstanceSpec(from *sqladmin.DatabaseInstance) *DatabaseInstanceSpec {
	if from == nil {
		return nil
	}
	return &DatabaseInstanceSpec{
		DatabaseVersion:      from.DatabaseVersion,
		InstanceType:         from.InstanceType,
		MasterInstanceName:   from.MasterInstanceName,
		Region:               from.Region,
		ReplicaConfiguration: newReplicaConfiguration(from.ReplicaConfiguration),
		Settings:             newSettings(from.Settings),
	}
}

func saveDatabaseInstanceSpec(from *DatabaseInstanceSpec) *sqladmin.DatabaseInstance {
	if from == nil {
		return nil
	}
	return &sqladmin.DatabaseInstance{
		DatabaseVersion:      from.DatabaseVersion,
		InstanceType:         from.InstanceType,
		MasterInstanceName:   from.MasterInstanceName,
		Region:               from.Region,
		ReplicaConfiguration: saveReplicaConfiguration(from.ReplicaConfiguration),
		Settings:             saveSettings(from.Settings),
	}
}

// DatabaseInstanceStatus: A Cloud SQL instance resource output.
type DatabaseInstanceStatus struct {
	// BackendType: FIRST_GEN: First Generation instance. MySQL only.
	// SECOND_GEN: Second Generation instance or PostgreSQL instance.
	// EXTERNAL: A database server that is not managed by Google.
	// This property is read-only; use the tier property in the Settings
	// object to determine the database type and Second or First Generation.
	BackendType string `json:"backendType,omitempty"`

	// ConnectionName: Connection name of the Cloud SQL instance used in connection strings.
	ConnectionName string `json:"connectionName,omitempty"`

	// FailoverReplica: The name and status of the failover replica. This
	// property is applicable only to Second Generation instances.
	FailoverReplica *DatabaseInstanceFailoverReplica `json:"failoverReplica,omitempty"`

	// GceZone: The Compute Engine zone that the instance is currently
	// serving from. This value could be different from the zone that was
	// specified when the instance was created if the instance has failed
	// over to its secondary zone.
	GceZone string `json:"gceZone,omitempty"`

	// IpAddresses: The assigned IP addresses for the instance.
	IpAddresses []*IpMapping `json:"ipAddresses,omitempty"`

	// Ipv6Address: The IPv6 address assigned to the instance. This property
	// is applicable only to First Generation instances.
	Ipv6Address string `json:"ipv6Address,omitempty"`

	// MaxDiskSize: The maximum disk size of the instance in bytes.
	MaxDiskSize int64 `json:"maxDiskSize,omitempty"`

	// Name: Name of the Cloud SQL instance. This does not include the
	// project ID.
	Name string `json:"name,omitempty"`

	// Project: The project ID of the project containing the Cloud SQL
	// instance. The Google apps domain is prefixed if applicable.
	Project string `json:"project,omitempty"`

	// ReplicaNames: The replicas of the instance.
	ReplicaNames []string `json:"replicaNames,omitempty"`

	// ServiceAccountEmailAddress: The service account email address
	// assigned to the instance. This property is applicable only to Second
	// Generation instances.
	ServiceAccountEmailAddress string `json:"serviceAccountEmailAddress,omitempty"`

	// State: The current serving state of the Cloud SQL instance. This can
	// be one of the following.
	// RUNNABLE: The instance is running, or is ready to run when accessed.
	// SUSPENDED: The instance is not available, for example due to problems with billing.
	// PENDING_CREATE: The instance is being created.
	// MAINTENANCE: The instance is down for maintenance.
	// FAILED: The instance creation failed.
	// UNKNOWN_STATE: The state of the instance is unknown.
	State string `json:"state,omitempty"`

	// SuspensionReason: If the instance state is SUSPENDED, the reason for the suspension.
	SuspensionReason []string `json:"suspensionReason,omitempty"`
}

func newDatabaseInstanceStatus(from *sqladmin.DatabaseInstance) *DatabaseInstanceStatus {
	var ips []*IpMapping
	for _, v := range from.IpAddresses {
		ips = append(ips, newIpMapping(v))
	}
	var repNames []string
	copy(repNames, from.ReplicaNames)

	return &DatabaseInstanceStatus{
		BackendType:                from.BackendType,
		ConnectionName:             from.ConnectionName,
		FailoverReplica:            newDatabaseInstanceFailoverReplica(from.FailoverReplica),
		GceZone:                    from.GceZone,
		IpAddresses:                ips,
		Ipv6Address:                from.Ipv6Address,
		MaxDiskSize:                from.MaxDiskSize,
		Name:                       from.Name,
		Project:                    from.Project,
		ReplicaNames:               repNames,
		ServiceAccountEmailAddress: from.ServiceAccountEmailAddress,
		State:                      from.State,
		SuspensionReason:           from.SuspensionReason,
	}
}

// IpConfiguration: IP Management configuration.
type IpConfiguration struct {
	// AuthorizedNetworks: The list of external networks that are allowed to
	// connect to the instance using the IP. In CIDR notation, also known as
	// 'slash' notation (e.g. 192.168.100.0/24).
	AuthorizedNetworks []*AclEntry `json:"authorizedNetworks,omitempty"`

	// Ipv4Enabled: Whether the instance should be assigned an IP address or
	// not.
	Ipv4Enabled bool `json:"ipv4Enabled,omitempty"`

	// PrivateNetwork: The resource link for the VPC network from which the
	// Cloud SQL instance is accessible for private IP. For example,
	// /projects/myProject/global/networks/default. This setting can be
	// updated, but it cannot be removed after it is set.
	PrivateNetwork string `json:"privateNetwork,omitempty"`

	// RequireSsl: Whether SSL connections over IP should be enforced or
	// not.
	RequireSsl bool `json:"requireSsl,omitempty"`
}

func newIpConfiguration(from *sqladmin.IpConfiguration) *IpConfiguration {
	if from == nil {
		return nil
	}
	var an []*AclEntry
	for _, v := range from.AuthorizedNetworks {
		an = append(an, newAclEntry(v))
	}
	return &IpConfiguration{
		AuthorizedNetworks: an,
		Ipv4Enabled:        from.Ipv4Enabled,
		PrivateNetwork:     from.PrivateNetwork,
		RequireSsl:         from.RequireSsl,
	}
}

func saveIpConfiguration(from *IpConfiguration) *sqladmin.IpConfiguration {
	if from == nil {
		return nil
	}

	var an []*sqladmin.AclEntry
	for _, v := range from.AuthorizedNetworks {
		an = append(an, saveAclEntry(v))
	}
	return &sqladmin.IpConfiguration{
		AuthorizedNetworks: an,
		Ipv4Enabled:        from.Ipv4Enabled,
		PrivateNetwork:     from.PrivateNetwork,
		RequireSsl:         from.RequireSsl,
		ForceSendFields:    []string{"Ipv4Enabled"},
	}
}

// IpMapping: Database instance IP Mapping.
type IpMapping struct {
	// IpAddress: The IP address assigned.
	IpAddress string `json:"ipAddress,omitempty"`

	// TimeToRetire: The due time for this IP to be retired in RFC 3339
	// format, for example 2012-11-15T16:19:00.094Z. This field is only
	// available when the IP is scheduled to be retired.
	TimeToRetire string `json:"timeToRetire,omitempty"`

	// Type: The type of this IP address. A PRIMARY address is an address
	// that can accept incoming connections. An OUTGOING address is the
	// source address of connections originating from the instance, if
	// supported.
	Type string `json:"type,omitempty"`
}

func newIpMapping(from *sqladmin.IpMapping) *IpMapping {
	if from == nil {
		return nil
	}
	return &IpMapping{
		IpAddress:    from.IpAddress,
		TimeToRetire: from.TimeToRetire,
		Type:         from.Type,
	}
}

func saveIpMapping(from *IpMapping) *sqladmin.IpMapping {
	if from == nil {
		return nil
	}
	return &sqladmin.IpMapping{
		IpAddress:    from.IpAddress,
		TimeToRetire: from.TimeToRetire,
		Type:         from.Type,
	}
}

// LocationPreference: Preferred location. This specifies where a Cloud
// SQL instance should preferably be located, either in a specific
// Compute Engine zone, or co-located with an App Engine application.
// Note that if the preferred location is not available, the instance
// will be located as close as possible within the region. Only one
// location may be specified.
type LocationPreference struct {
	// FollowGaeApplication: The AppEngine application to follow, it must be
	// in the same region as the Cloud SQL instance.
	FollowGaeApplication string `json:"followGaeApplication,omitempty"`

	// Zone: The preferred Compute Engine zone (e.g. us-central1-a,
	// us-central1-b, etc.).
	Zone string `json:"zone,omitempty"`
}

func newLocationPreference(from *sqladmin.LocationPreference) *LocationPreference {
	if from == nil {
		return nil
	}
	return &LocationPreference{
		FollowGaeApplication: from.FollowGaeApplication,
		Zone:                 from.Zone,
	}
}

func saveLocationPreference(from *LocationPreference) *sqladmin.LocationPreference {
	if from == nil {
		return nil
	}
	return &sqladmin.LocationPreference{
		FollowGaeApplication: from.FollowGaeApplication,
		Zone:                 from.Zone,
	}
}

// MaintenanceWindow: Maintenance window. This specifies when a v2 Cloud
// SQL instance should preferably be restarted for system maintenance
// purposes.
type MaintenanceWindow struct {
	// Day: day of week (1-7), starting on Monday.
	Day int64 `json:"day,omitempty"`

	// Hour: hour of day - 0 to 23.
	Hour int64 `json:"hour,omitempty"`

	// UpdateTrack: Maintenance timing setting: canary (Earlier) or stable
	// (Later).
	//  Learn more.
	UpdateTrack string `json:"updateTrack,omitempty"`
}

func newMainenanceWindow(from *sqladmin.MaintenanceWindow) *MaintenanceWindow {
	if from == nil {
		return nil
	}
	return &MaintenanceWindow{
		Day:         from.Day,
		Hour:        from.Hour,
		UpdateTrack: from.UpdateTrack,
	}
}

func saveMaintenanceWindow(from *MaintenanceWindow) *sqladmin.MaintenanceWindow {
	if from == nil {
		return nil
	}
	return &sqladmin.MaintenanceWindow{
		Day:         from.Day,
		Hour:        from.Hour,
		UpdateTrack: from.UpdateTrack,
	}
}

// MySqlReplicaConfiguration: Read-replica configuration specific to
// MySQL databases.
type MySqlReplicaConfiguration struct {
	// CaCertificate: PEM representation of the trusted CA's x509 certificate.
	CaCertificate string `json:"caCertificate,omitempty"`

	// ClientCertificate: PEM representation of the slave's x509 certificate.
	ClientCertificate string `json:"clientCertificate,omitempty"`

	// ClientKey: PEM representation of the slave's private key. The
	// corresponsing public key is encoded in the client's certificate.
	ClientKey string `json:"clientKey,omitempty"`

	// ConnectRetryInterval: Seconds to wait between connect retries.
	// MySQL's default is 60 seconds.
	ConnectRetryInterval int64 `json:"connectRetryInterval,omitempty"`

	// DumpFilePath: Path to a SQL dump file in Google Cloud Storage from
	// which the slave instance is to be created. The URI is in the form
	// gs://bucketName/fileName. Compressed gzip files (.gz) are also
	// supported. Dumps should have the binlog co-ordinates from which
	// replication should begin. This can be accomplished by setting
	// --master-data to 1 when using mysqldump.
	DumpFilePath string `json:"dumpFilePath,omitempty"`

	// MasterHeartbeatPeriod: Interval in milliseconds between replication heartbeats.
	MasterHeartbeatPeriod int64 `json:"masterHeartbeatPeriod,omitempty"`

	// Password: The password for the replication connection.
	Password string `json:"password,omitempty"`

	// SslCipher: A list of permissible ciphers to use for SSL encryption.
	SslCipher string `json:"sslCipher,omitempty"`

	// Username: The username for the replication connection.
	Username string `json:"username,omitempty"`

	// VerifyServerCertificate: Whether or not to check the master's Common
	// Name value in the certificate that it sends during the SSL handshake.
	VerifyServerCertificate bool `json:"verifyServerCertificate,omitempty"`
}

func newMySqlReplicationConfiguration(from *sqladmin.MySqlReplicaConfiguration) *MySqlReplicaConfiguration {
	if from == nil {
		return nil
	}
	return &MySqlReplicaConfiguration{
		CaCertificate:           from.CaCertificate,
		ClientCertificate:       from.ClientCertificate,
		ClientKey:               from.ClientKey,
		ConnectRetryInterval:    from.ConnectRetryInterval,
		DumpFilePath:            from.DumpFilePath,
		MasterHeartbeatPeriod:   from.MasterHeartbeatPeriod,
		Password:                from.Password,
		SslCipher:               from.SslCipher,
		Username:                from.Username,
		VerifyServerCertificate: from.VerifyServerCertificate,
	}
}

func saveMySqlReplicaConfiguration(from *MySqlReplicaConfiguration) *sqladmin.MySqlReplicaConfiguration {
	if from == nil {
		return nil
	}
	return &sqladmin.MySqlReplicaConfiguration{
		CaCertificate:           from.CaCertificate,
		ClientCertificate:       from.ClientCertificate,
		ClientKey:               from.ClientKey,
		ConnectRetryInterval:    from.ConnectRetryInterval,
		DumpFilePath:            from.DumpFilePath,
		MasterHeartbeatPeriod:   from.MasterHeartbeatPeriod,
		Password:                from.Password,
		SslCipher:               from.SslCipher,
		Username:                from.Username,
		VerifyServerCertificate: from.VerifyServerCertificate,
	}
}

// ReplicaConfiguration: Read-replica configuration for connecting to
// the master.
type ReplicaConfiguration struct {
	// FailoverTarget: Specifies if the replica is the failover target. If
	// the field is set to true the replica will be designated as a failover
	// replica. In case the master instance fails, the replica instance will
	// be promoted as the new master instance.
	// Only one replica can be specified as failover target, and the replica
	// has to be in different zone with the master instance.
	FailoverTarget bool `json:"failoverTarget,omitempty"`

	// MysqlReplicaConfiguration: MySQL specific configuration when
	// replicating from a MySQL on-premises master. Replication
	// configuration information such as the username, password,
	// certificates, and keys are not stored in the instance metadata. The
	// configuration information is used only to set up the replication
	// connection and is stored by MySQL in a file named master.info in the
	// data directory.
	MysqlReplicaConfiguration *MySqlReplicaConfiguration `json:"mysqlReplicaConfiguration,omitempty"`
}

func newReplicaConfiguration(from *sqladmin.ReplicaConfiguration) *ReplicaConfiguration {
	if from == nil {
		return nil
	}
	return &ReplicaConfiguration{
		FailoverTarget:            from.FailoverTarget,
		MysqlReplicaConfiguration: newMySqlReplicationConfiguration(from.MysqlReplicaConfiguration),
	}
}

func saveReplicaConfiguration(from *ReplicaConfiguration) *sqladmin.ReplicaConfiguration {
	if from == nil {
		return nil
	}
	return &sqladmin.ReplicaConfiguration{
		FailoverTarget: from.FailoverTarget,
		//MysqlReplicaConfiguration: newMySqlReplicationConfiguration(from.MysqlReplicaConfiguration),
	}
}

// Settings: Database instance Settings.
type Settings struct {
	// ActivationPolicy: The activation policy specifies when the instance
	// is activated; it is applicable only when the instance state is
	// RUNNABLE. Valid values:
	// ALWAYS: The instance is on, and remains so even in the absence of
	// connection requests.
	// NEVER: The instance is off; it is not activated, even if a connection
	// request arrives.
	// ON_DEMAND: First Generation instances only. The instance responds to
	// incoming requests, and turns itself off when not in use. Instances
	// with PER_USE pricing turn off after 15 minutes of inactivity.
	// Instances with PER_PACKAGE pricing turn off after 12 hours of
	// inactivity.
	ActivationPolicy string `json:"activationPolicy,omitempty"`

	// AuthorizedGaeApplications: The App Engine app IDs that can access this instance. First Generation instances only.
	AuthorizedGaeApplications []string `json:"authorizedGaeApplications,omitempty"`

	// AvailabilityType: Availability type (PostgreSQL instances only).
	// Potential values:
	// ZONAL: The instance serves data from only one zone. Outages in that zone affect data accessibility.
	// REGIONAL: The instance can serve data from more than one zone in a region (it is highly available).
	// For more information, see Overview of the High Availability Configuration.
	AvailabilityType string `json:"availabilityType,omitempty"`

	// BackupConfiguration: The daily backup configuration for the instance.
	BackupConfiguration *BackupConfiguration `json:"backupConfiguration,omitempty"`

	// CrashSafeReplicationEnabled: Configuration specific to read replica
	// instances. Indicates whether database flags for crash-safe
	// replication are enabled. This property is only applicable to First
	// Generation instances.
	CrashSafeReplicationEnabled bool `json:"crashSafeReplicationEnabled,omitempty"`

	// DataDiskSizeGb: The size of data disk, in GB. The data disk size
	// minimum is 10GB. Not used for First Generation instances.
	DataDiskSizeGb int64 `json:"dataDiskSizeGb,omitempty"`

	// DataDiskType: The type of data disk: PD_SSD (default) or PD_HDD. Not
	// used for First Generation instances.
	DataDiskType string `json:"dataDiskType,omitempty"`

	// DatabaseFlags: The database flags passed to the instance at startup.
	DatabaseFlags []*DatabaseFlags `json:"databaseFlags,omitempty"`

	// DatabaseReplicationEnabled: Configuration specific to read replica
	// instances. Indicates whether replication is enabled or not.
	DatabaseReplicationEnabled bool `json:"databaseReplicationEnabled,omitempty"`

	// IpConfiguration: The Settings for IP Management. This allows to
	// enable or disable the instance IP and manage which external networks
	// can connect to the instance. The IPv4 address cannot be disabled for
	// Second Generation instances.
	IpConfiguration *IpConfiguration `json:"ipConfiguration,omitempty"`

	// LocationPreference: The location preference Settings. This allows the
	// instance to be located as near as possible to either an App Engine
	// app or Compute Engine zone for better performance. App Engine
	// co-location is only applicable to First Generation instances.
	LocationPreference *LocationPreference `json:"locationPreference,omitempty"`

	// MaintenanceWindow: The maintenance window for this instance. This
	// specifies when the instance can be restarted for maintenance
	// purposes. Not used for First Generation instances.
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// PricingPlan: The pricing plan for this instance. This can be either
	// PER_USE or PACKAGE. Only PER_USE is supported for Second Generation
	// instances.
	PricingPlan string `json:"pricingPlan,omitempty"`

	// ReplicationType: The type of replication this instance uses. This can
	// be either ASYNCHRONOUS or SYNCHRONOUS. This property is only
	// applicable to First Generation instances.
	ReplicationType string `json:"replicationType,omitempty"`

	// SettingsVersion: The version of instance Settings. This is a required
	// field for update method to make sure concurrent updates are handled
	// properly. During update, use the most recent settingsVersion value
	// for this instance and do not try to update this value.
	// TODO: do we need this here?
	SettingsVersion int64 `json:"settingsVersion,omitempty"`

	// StorageAutoResize: Configuration to increase storage size
	// automatically. The default value is true. Not used for First
	// Generation instances.
	StorageAutoResize *bool `json:"storageAutoResize,omitempty"`

	// StorageAutoResizeLimit: The maximum size to which storage capacity
	// can be automatically increased. The default value is 0, which
	// specifies that there is no limit. Not used for First Generation
	// instances.
	StorageAutoResizeLimit int64 `json:"storageAutoResizeLimit,omitempty"`

	// Tier: The tier (or machine type) for this instance, for example
	// db-n1-standard-1 (MySQL instances) or db-custom-1-3840 (PostgreSQL
	// instances). For MySQL instances, this property determines whether the
	// instance is First or Second Generation. For more information, see
	// Instance Settings.
	Tier string `json:"tier,omitempty"`

	// UserLabels: User-provided labels, represented as a dictionary where
	// each label is a single key value pair.
	UserLabels map[string]string `json:"userLabels,omitempty"`
}

func newSettings(from *sqladmin.Settings) *Settings {
	if from == nil {
		return nil
	}
	var df []*DatabaseFlags
	for _, v := range from.DatabaseFlags {
		df = append(df, newDatabaseFlags(v))
	}
	return &Settings{
		ActivationPolicy:            from.ActivationPolicy,
		AuthorizedGaeApplications:   from.AuthorizedGaeApplications,
		AvailabilityType:            from.AvailabilityType,
		BackupConfiguration:         newBackupConfiguration(from.BackupConfiguration),
		CrashSafeReplicationEnabled: from.CrashSafeReplicationEnabled,
		DataDiskSizeGb:              from.DataDiskSizeGb,
		DataDiskType:                from.DataDiskType,
		DatabaseFlags:               df,
		DatabaseReplicationEnabled:  from.DatabaseReplicationEnabled,
		IpConfiguration:             newIpConfiguration(from.IpConfiguration),
		LocationPreference:          newLocationPreference(from.LocationPreference),
		MaintenanceWindow:           newMainenanceWindow(from.MaintenanceWindow),
		PricingPlan:                 from.PricingPlan,
		ReplicationType:             from.ReplicationType,
		StorageAutoResize:           from.StorageAutoResize,
		StorageAutoResizeLimit:      from.StorageAutoResizeLimit,
		Tier:                        from.Tier,
		UserLabels:                  from.UserLabels,
	}
}

func saveSettings(from *Settings) *sqladmin.Settings {
	if from == nil {
		return nil
	}
	var df []*sqladmin.DatabaseFlags
	for _, v := range from.DatabaseFlags {
		df = append(df, saveDatabaseFlags(v))
	}
	return &sqladmin.Settings{
		ActivationPolicy:            from.ActivationPolicy,
		AuthorizedGaeApplications:   from.AuthorizedGaeApplications,
		AvailabilityType:            from.AvailabilityType,
		BackupConfiguration:         saveBackupConfiguration(from.BackupConfiguration),
		CrashSafeReplicationEnabled: from.CrashSafeReplicationEnabled,
		DataDiskSizeGb:              from.DataDiskSizeGb,
		DataDiskType:                from.DataDiskType,
		DatabaseFlags:               df,
		DatabaseReplicationEnabled:  from.DatabaseReplicationEnabled,
		IpConfiguration:             saveIpConfiguration(from.IpConfiguration),
		LocationPreference:          saveLocationPreference(from.LocationPreference),
		MaintenanceWindow:           saveMaintenanceWindow(from.MaintenanceWindow),
		PricingPlan:                 from.PricingPlan,
		ReplicationType:             from.ReplicationType,
		StorageAutoResize:           from.StorageAutoResize,
		StorageAutoResizeLimit:      from.StorageAutoResizeLimit,
		Tier:                        from.Tier,
		UserLabels:                  from.UserLabels,
	}
}
