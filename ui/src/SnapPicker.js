import React from "react";
import { useHistory } from "react-router-dom";

import { Select, MenuItem } from "@material-ui/core";

import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en";

TimeAgo.addLocale(en);

/**
 * Picker is the actual MaterialUI component for SnapPicker
 * @param {Object} props - React properties
 * @param {Object[]} props.items - List of restic snapshots (`/api/v1/snapshots`)
 * @param {string} props.current - `snapshot.id` of the selected snapshot
 * @param {Function} props.onChange - `onChange` handler
 */
const Picker = ({ items, current, onChange }) => {
  const timeAgo = new TimeAgo("en-US");
  return (
    <Select disabled={items.length === 0} onChange={onChange} value={current}>
      {items.map((s) => (
        <MenuItem key={s.id} value={s.id}>
          {timeAgo.format(new Date(s.time), "twitter")}
        </MenuItem>
      ))}
    </Select>
  );
};

/**
 * SnapPicker is the upper-right `<select>` component for choosing the snapshot.
 * It currently queries it's own data and displays a snackbar if that fails.
 */
const SnapPicker = ({ snapshots, job, path, current }) => {
  const data = snapshots
    .map((s) => {
      s.id = s.id.substring(0, 8);
      return s;
    })
    .reverse();

  const history = useHistory();

  if (current === "latest") {
    current = data[0]?.id;
  }

  return (
    <Picker
      current={data.length > 0 ? current : ""}
      items={data}
      onChange={(e) => {
        let dest = `/${job}/${e.target.value}`;
        if (path) {
          dest += `/${path}`;
        }
        history.push(dest);
      }}
    />
  );
};

export default SnapPicker;
