import React, { useMemo, useState } from "react";
import { useHistory } from "react-router-dom";

import {
  Button,
  Divider,
  IconButton,
  ListItemText,
  Menu,
  MenuItem,
  Popover,
  Tooltip,
  Typography,
} from "@material-ui/core";
import { useTheme } from "@material-ui/core/styles";
import { CalendarToday, ChevronLeft, ChevronRight } from "@material-ui/icons";

import {
  formatBytes,
  formatDateTime,
  formatInteger,
} from "./format";

const WEEKDAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

/**
 * SnapPicker is the upper-right calendar component for choosing the snapshot.
 */
const SnapPicker = ({ snapshots, job, path, current }) => {
  const data = useMemo(() => normalizeSnapshots(snapshots), [snapshots]);
  const history = useHistory();
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState(null);
  const selected = selectedSnapshot(data, current);
  const [visibleMonth, setVisibleMonth] = useState(() =>
    monthStart(selected ? new Date(selected.time) : new Date())
  );

  React.useEffect(() => {
    if (selected) {
      setVisibleMonth(monthStart(new Date(selected.time)));
    }
  }, [selected && selected.id]);

  const snapshotsByDay = useMemo(() => groupSnapshotsByDay(data), [data]);

  const navigate = (snapshot) => {
    let dest = `/${job}/${snapshot.shortId}`;
    if (path) {
      dest += `/${path}`;
    }
    history.push(dest);
    setAnchorEl(null);
  };

  return (
    <>
      <Button
        color="inherit"
        disabled={data.length === 0}
        onClick={(event) => setAnchorEl(event.currentTarget)}
      >
        <CalendarToday fontSize="small" style={{ marginRight: 6 }} />
        {selected ? formatButtonDate(selected.time) : "No snapshots"}
      </Button>
      <Popover
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
        transformOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <div
          style={{
            width: 322,
            padding: "0.75em",
            fontFamily: theme.typography.fontFamily,
          }}
        >
          <CalendarHeader
            visibleMonth={visibleMonth}
            onPrevious={() => setVisibleMonth(addMonths(visibleMonth, -1))}
            onNext={() => setVisibleMonth(addMonths(visibleMonth, 1))}
          />
          <Divider style={{ marginBottom: "0.5em" }} />
          <CalendarGrid
            visibleMonth={visibleMonth}
            selected={selected}
            snapshotsByDay={snapshotsByDay}
            onSelect={navigate}
          />
        </div>
      </Popover>
    </>
  );
};

const CalendarHeader = ({ visibleMonth, onPrevious, onNext }) => (
  <div
    style={{
      display: "flex",
      alignItems: "center",
      justifyContent: "space-between",
      marginBottom: "0.5em",
    }}
  >
    <IconButton size="small" onClick={onPrevious} aria-label="Previous month">
      <ChevronLeft />
    </IconButton>
    <Typography variant="subtitle1">
      {visibleMonth.toLocaleDateString(undefined, {
        month: "long",
        year: "numeric",
      })}
    </Typography>
    <IconButton size="small" onClick={onNext} aria-label="Next month">
      <ChevronRight />
    </IconButton>
  </div>
);

const CalendarGrid = ({ visibleMonth, selected, snapshotsByDay, onSelect }) => {
  const days = calendarDays(visibleMonth);
  const selectedDay = selected ? dayKey(new Date(selected.time)) : undefined;
  const currentMonth = visibleMonth.getMonth();
  const [snapshotMenuAnchor, setSnapshotMenuAnchor] = useState(null);
  const [snapshotMenuItems, setSnapshotMenuItems] = useState([]);

  const handleDayClick = (event, daySnapshots) => {
    if (daySnapshots.length === 0) return;

    if (daySnapshots.length === 1) {
      onSelect(daySnapshots[0]);
      return;
    }

    setSnapshotMenuAnchor(event.currentTarget);
    setSnapshotMenuItems([...daySnapshots].reverse());
  };

  const handleSnapshotSelect = (snapshot) => {
    setSnapshotMenuAnchor(null);
    setSnapshotMenuItems([]);
    onSelect(snapshot);
  };

  return (
    <div>
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(7, 1fr)",
          gap: 4,
          marginBottom: 4,
        }}
      >
        {WEEKDAYS.map((day) => (
          <Typography
            key={day}
            variant="caption"
            color="textSecondary"
            align="center"
          >
            {day}
          </Typography>
        ))}
      </div>
      <div
        style={{ display: "grid", gridTemplateColumns: "repeat(7, 1fr)", gap: 4 }}
      >
        {days.map((date) => {
          const key = dayKey(date);
          const daySnapshots = snapshotsByDay[key] || [];
          const hasSnapshot = daySnapshots.length > 0;
          const isSelected = key === selectedDay;
          const isOtherMonth = date.getMonth() !== currentMonth;

          return (
            <Tooltip key={key} title={tooltipFor(date, daySnapshots)}>
              <div
                onClick={(event) => handleDayClick(event, daySnapshots)}
                style={{
                  position: "relative",
                  height: 36,
                  lineHeight: "36px",
                  textAlign: "center",
                  borderRadius: 18,
                  cursor: hasSnapshot ? "pointer" : "default",
                  color: isSelected
                    ? "white"
                    : hasSnapshot
                    ? "inherit"
                    : "#9e9e9e",
                  background: isSelected
                    ? "#3f51b5"
                    : hasSnapshot
                    ? "rgba(63, 81, 181, 0.10)"
                    : "transparent",
                  opacity: isOtherMonth ? 0.45 : 1,
                  fontWeight: hasSnapshot ? 600 : 400,
                  userSelect: "none",
                  fontFamily: "inherit",
                }}
              >
                {date.getDate()}
                {daySnapshots.length > 1 && (
                  <span
                    style={{
                      position: "absolute",
                      right: 5,
                      bottom: 1,
                      fontSize: 9,
                      lineHeight: "9px",
                    }}
                  >
                    {daySnapshots.length}
                  </span>
                )}
              </div>
            </Tooltip>
          );
        })}
      </div>
      <Menu
        anchorEl={snapshotMenuAnchor}
        open={Boolean(snapshotMenuAnchor)}
        onClose={() => setSnapshotMenuAnchor(null)}
      >
        {snapshotMenuItems.map((snapshot) => {
          const summary = snapshot.summary || {};
          return (
            <MenuItem
              key={snapshot.id}
              selected={selected && selected.shortId === snapshot.shortId}
              onClick={() => handleSnapshotSelect(snapshot)}
            >
              <ListItemText
                primary={formatButtonDate(snapshot.time)}
                secondary={`Size ${formatBytes(
                  summary.total_bytes_processed
                )} · files ${formatInteger(summary.total_files_processed)}`}
              />
            </MenuItem>
          );
        })}
      </Menu>
      <Typography
        variant="caption"
        color="textSecondary"
        style={{ display: "block", marginTop: "0.75em" }}
      >
        Grey days have no snapshot. Hover a day for details; if multiple snapshots
        exist on a day, clicking opens a chooser.
      </Typography>
    </div>
  );
};

function normalizeSnapshots(snapshots) {
  return snapshots
    .map((snapshot) => ({
      ...snapshot,
      shortId: snapshot.id.substring(0, 8),
    }))
    .sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
}

function selectedSnapshot(snapshots, current) {
  if (snapshots.length === 0) return undefined;
  if (current === "latest") return snapshots[snapshots.length - 1];

  return snapshots.find(
    (snapshot) => snapshot.id === current || snapshot.shortId === current
  );
}

function groupSnapshotsByDay(snapshots) {
  return snapshots.reduce((groups, snapshot) => {
    const key = dayKey(new Date(snapshot.time));
    groups[key] = [...(groups[key] || []), snapshot];
    return groups;
  }, {});
}

function calendarDays(visibleMonth) {
  const first = monthStart(visibleMonth);
  const start = new Date(first);
  start.setDate(first.getDate() - first.getDay());

  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(start);
    date.setDate(start.getDate() + index);
    return date;
  });
}

function monthStart(date) {
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function addMonths(date, count) {
  return new Date(date.getFullYear(), date.getMonth() + count, 1);
}

function dayKey(date) {
  return [date.getFullYear(), date.getMonth() + 1, date.getDate()]
    .map((part) => String(part).padStart(2, "0"))
    .join("-");
}

function formatButtonDate(value) {
  return new Date(value).toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function tooltipFor(date, snapshots) {
  if (snapshots.length === 0) {
    return `No snapshot on ${date.toLocaleDateString()}`;
  }

  const latest = snapshots[snapshots.length - 1];
  const summary = latest.summary || {};
  const files = summary.total_files_processed;
  const size = summary.total_bytes_processed;

  return `${snapshots.length} snapshot${snapshots.length === 1 ? "" : "s"}; latest ${formatDateTime(
    latest.time
  )}; size ${formatBytes(size)}; files ${formatInteger(files)}`;
}

export default SnapPicker;
