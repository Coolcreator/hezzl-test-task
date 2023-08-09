CREATE TABLE IF NOT EXISTS hezzl.logs (
    Id UInt64,
    ProjectId UInt64,
    Name String,
    Description String,
    Priority UInt32,
    Removed Boolean DEFAULT false,
    EventTime DateTime
) ENGINE = MergeTree()
ORDER BY (Id, ProjectId, Name);

ALTER TABLE hezzl.logs ADD INDEX idx_Id(Id) TYPE minmax GRANULARITY 1;

ALTER TABLE hezzl.logs ADD INDEX idx_ProjectId(ProjectId) TYPE minmax GRANULARITY 1;

ALTER TABLE hezzl.logs ADD INDEX idx_Name(Name) TYPE bloom_filter GRANULARITY 1;
