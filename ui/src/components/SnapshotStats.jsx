import React from "react";
import {
  Divider,
  Grid,
  List,
  ListItem,
  ListItemText,
  Typography,
} from "@material-ui/core";

import {
  formatBytes,
  formatDateTime,
  formatDuration,
  formatInteger,
} from "../format";

/**
 * SnapshotStats renders the selected restic snapshot's metadata below the file
 * list, when the repository provides restic summary data.
 */
export const SnapshotStats = ({ snapshot }) => {
  if (!snapshot) return null;

  const summary = snapshot.summary || {};
  const fileCount =
    summary.total_files_processed !== undefined && summary.total_files_processed !== null
      ? summary.total_files_processed
      : sumKnown(summary.files_new, summary.files_changed, summary.files_unmodified);
  const dirCount = sumKnown(
    summary.dirs_new,
    summary.dirs_changed,
    summary.dirs_unmodified
  );

  const metrics = [
    ["Snapshot time", formatDateTime(snapshot.time)],
    ["Snapshot size", formatBytes(summary.total_bytes_processed)],
    ["Files", formatInteger(fileCount)],
    ["Directories", formatInteger(dirCount)],
    ["Data added", formatBytes(summary.data_added)],
    ["Packed data added", formatBytes(summary.data_added_packed)],
    ["Backup duration", formatDuration(summary.backup_start, summary.backup_end)],
    ["Hostname", snapshot.hostname || "—"],
  ];

  return (
    <>
      <Divider />
      <div style={{ padding: "1em" }}>
        <Typography variant="subtitle2" color="textSecondary" gutterBottom>
          Snapshot details
        </Typography>
        <Grid container spacing={1}>
          {metrics.map(([label, value]) => (
            <Grid item xs={12} sm={6} md={3} key={label}>
              <Typography variant="caption" color="textSecondary">
                {label}
              </Typography>
              <Typography variant="body2">{value}</Typography>
            </Grid>
          ))}
        </Grid>

        {snapshot.summary && (
          <List dense style={{ marginTop: "0.5em" }}>
            <ListItem disableGutters>
              <ListItemText
                primary="Changed in this backup"
                secondary={`${formatInteger(summary.files_new)} new files, ${formatInteger(
                  summary.files_changed
                )} changed files, ${formatInteger(
                  summary.files_unmodified
                )} unchanged files`}
              />
            </ListItem>
          </List>
        )}
      </div>
    </>
  );
};

function sumKnown(...values) {
  const known = values.filter((value) => value !== undefined && value !== null);
  if (known.length === 0) return undefined;

  return known.reduce((sum, value) => sum + Number(value), 0);
}
