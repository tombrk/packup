/** @jsx jsx */
import { jsx } from "@emotion/core";
import { Component } from "react";
import Axios from "axios";
import queryString from "query-string";
import { formatRelative } from "date-fns";

import { Select, MenuItem } from "@material-ui/core";

const api =
  process.env.NODE_ENV === "development"
    ? "http://localhost:2112/api/v1"
    : "/api/v1";

export default class SnapshotPicker extends Component {
  state = {
    error: null,
    loaded: false,
    items: [],
  };

  componentDidMount() {
    Axios.get(`${api}/snapshots`).then(
      (result) =>
        this.setState({
          loaded: true,
          items: result.data == null ? [] : result.data,
        }),
      (error) => {
        this.setState({ loaded: true, error: error, items: [] });
      }
    );
  }

  render() {
    const { items } = this.state;
    const qv = queryString.parse(this.props.location.search);

    return (
      <Select
        disabled={items.length === 0}
        value={qv.snapshot === undefined ? "latest" : qv.snapshot}
        onChange={(e) => {
          this.props.history.push(
            `${this.props.location.pathname}?${queryString.stringify({
              ...qv,
              snapshot: e.target.value,
            })}`
          );
        }}
      >
        <MenuItem value="latest">latest</MenuItem>
        {items.reverse().map((item) => (
          <MenuItem key={item.id} value={item.id}>
            {formatRelative(new Date(item.time), new Date())}
          </MenuItem>
        ))}
      </Select>
    );
  }
}
