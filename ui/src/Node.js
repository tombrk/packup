/** @jsx jsx */
import { jsx } from "@emotion/core";

import { Link } from "react-router-dom";
import { IconButton } from "@material-ui/core";
import {
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemSecondaryAction
} from "@material-ui/core";
import { InsertDriveFile, Archive, Folder } from "@material-ui/icons";
import ContentLoader from "react-content-loader";

const api =
  process.env.NODE_ENV === "development"
    ? "http://localhost:2112/api/v1"
    : "/api/v1";

const ListItemLink = props => <ListItem button {...props} />;
const ALink = props => <Link {...props} to={props.href} />;

const Node = props => {
  const Icon = props.icon;
  return (
    <ListItemLink component={props.link} href={props.href}>
      <ListItemIcon>
        <Icon />
      </ListItemIcon>
      <ListItemText primary={props.name} {...props.text} />
      {props.children}
    </ListItemLink>
  );
};

export const File = props => (
  <Node
    {...props}
    link="a"
    href={`${api}/dump?path=${props.path}`}
    icon={InsertDriveFile}
    text={{ secondary: formatBytes(props.size) }}
  />
);

export const Dir = props => (
  <Node
    {...props}
    name={props.name + "/"}
    link={ALink}
    href={`${props.path}${window.location.search}`}
    icon={Folder}
  >
    <ListItemSecondaryAction>
      <IconButton
        component="a"
        href={`${api}/dump?path=${
          props.path
        }&compress=true&filename=${props.path.split(/[\\/]/).pop()}.tar.gz`}
      >
        <Archive />
      </IconButton>
    </ListItemSecondaryAction>
  </Node>
);

export const Placeholder = () => (
  <ListItem button>
    <ListItemIcon>
      <InsertDriveFile />
    </ListItemIcon>
    <ListItemText>
      <ContentLoader css={{ height: "1em", width: "100%" }} />
    </ListItemText>
  </ListItem>
);

function formatBytes(a, b) {
  if (0 === a) return "0 Bytes";
  var c = 1024,
    d = b || 2,
    e = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"],
    f = Math.floor(Math.log(a) / Math.log(c));
  return parseFloat((a / Math.pow(c, f)).toFixed(d)) + " " + e[f];
}
