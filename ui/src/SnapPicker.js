import React from "react";
import { useHistory, useParams } from "react-router-dom";

import { formatRelative } from "date-fns";

import { Select, MenuItem } from "@material-ui/core";

/**
 * Picker is the actual MaterialUI component for SnapPicker
 * @param {Object} props - React properties
 * @param {Object[]} props.items - List of restic snapshots (`/api/v1/snapshots`)
 * @param {string} props.current - `snapshot.id` of the selected snapshot
 * @param {Function} props.onChange - `onChange` handler
 */
const Picker = ({ items, current, onChange }) => {
  return (
    <Select disabled={items.length === 0} onChange={onChange} value={current}>
      {items.map((s) => (
        <MenuItem key={s.id} value={s.id}>
          {formatRelative(new Date(s.time), new Date())}
        </MenuItem>
      ))}
    </Select>
  );
};

/**
 * SnapPicker is the upper-right `<select>` component for choosing the snapshot.
 * It currently queries it's own data and displays a snackbar if that fails.
 */
const SnapPicker = ({ snapshots, job, path }) => {
  const data = snapshots
    .map((s) => {
      s.id = s.id.substring(0, 8);
      return s;
    })
    .reverse();

  const history = useHistory();

  return (
    <Picker
      current={data.length > 0 ? data[0].id : ""}
      items={data}
      onChange={(e) => {
        history.push(`/${job}/${e.target.value}/${path}`);
      }}
    />
  );
};

export default SnapPicker;
