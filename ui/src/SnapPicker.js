import React, { useState, useEffect } from "react";

import axios from "axios";
import queryString from "query-string";
import { formatRelative } from "date-fns";
import { addr } from "./api";

import { Select, MenuItem, Button } from "@material-ui/core";
import { useSnackbar } from "notistack";

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
        <MenuItem value={s.id} value={s.id}>
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
const SnapPicker = ({ history, location }) => {
  const [data, setData] = useState([]);
  const qv = queryString.parse(location.search);
  const { enqueueSnackbar } = useSnackbar();

  const push = (to) => {
    history.push(
      `${location.pathname}?${queryString.stringify({
        ...qv,
        snapshot: to,
      })}`
    );
  };

  useEffect(() => {
    const fetch = async () => {
      try {
        const result = await axios(`${addr}/snapshots`);
        setData(
          result.data.reverse().map((s) => {
            s.id = s.id.substring(0, 8);
            return s;
          })
        );
      } catch (error) {
        const message = error.response
          ? error.response.data
          : "Unable to reach backend. Please check your network connection";

        enqueueSnackbar(message, {
          variant: "error",
        });
      }
    };
    fetch();
  }, []);

  return (
    <Picker
      current={data.length > 0 ? data[0].id : ""}
      items={data}
      onChange={(e) => {
        push(e.target.value);
      }}
    ></Picker>
  );
};

export default SnapPicker;
