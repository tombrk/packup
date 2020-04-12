/** @jsx jsx */
import { jsx } from "@emotion/core";
import queryString from "query-string";

import { Link as RouterLink } from "react-router-dom";
import { Breadcrumbs, Link, Typography } from "@material-ui/core";
import { AppBar, Toolbar } from "@material-ui/core";

export const Path = props => {
  const qv = queryString.parse(window.location.search);
  if (qv.snapshot === undefined) {
    qv.snapshot = "latest";
  }

  const breads = props.dir
    .split("/")
    .reduce((total, current) => {
      if (props.dir === "/") {
        total = [""];
      } else {
        total.push(current);
      }
      return total;
    }, [])
    .map(node => (node === "" ? qv.snapshot : node))
    .map((node, index) => {
      const url =
        node === qv.snapshot
          ? "/"
          : props.dir
              .split("/")
              .slice(0, index + 1)
              .join("/");

      return (
        <Link
          component={RouterLink}
          key={url}
          css={{ marginRight: ".2em" }}
          to={`${url}${window.location.search}`}
        >
          {node}
        </Link>
      );
    });

  return <Breadcrumbs>{breads}</Breadcrumbs>;
};

export const TitleBar = props => (
  <AppBar position="static">
    <Toolbar css={{ display: "flex" }}>
      <Typography css={{ marginRight: "1em" }} variant="h6">
        prestic
      </Typography>
      {props.children}
    </Toolbar>
  </AppBar>
);
