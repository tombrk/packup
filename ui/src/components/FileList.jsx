import { Folder, InsertDriveFile, Archive } from "@material-ui/icons";
import { Link, useLocation, useParams } from "react-router-dom";
import {
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
} from "@material-ui/core";
import React from "react";

import { addr } from "../api";
import { formatBytes } from "../format";

/**
 * FileList lists files from a path of a Restic snapshot
 * @param {Object} props - React props
 * @param {string} props.files - List of `restic.Node`
 */
export const FileList = ({ files }) => {
  return (
    <List>
      {files.map((f) => {
        const Node = f.type === "dir" ? Dir : File;
        return <Node key={f.name} {...f}></Node>;
      })}
    </List>
  );
};

/**
 * Directory renders a `ListItem` for a directory
 * @param {Object} props - React props
 * @param {String} props.name - Name of the directory (no slashes)
 */
const Dir = ({ name, path }) => (
  <ItemIconLink
    name={`${name}/`}
    to={`${useLocation().pathname}/${name}`}
    icon={<Folder />}
  >
    <ListItemSecondaryAction>
      <DownloadButton path={path} />
    </ListItemSecondaryAction>
  </ItemIconLink>
);

/**
 * DownloadButton renders the Button on the side of a directory Node,
 * to download it entirely as a tar.gz
 * @param {Object} props - React props
 * @param {string} props.path - Directory path (from API)
 */
const DownloadButton = ({ path }) => {
  const { job, snapshot } = useParams();
  const to = `${addr}/jobs/${job}/dump?path=${path}&snapshot=${snapshot}&compress=true`;
  return (
    <IconButton component="a" href={to}>
      <Archive></Archive>
    </IconButton>
  );
};

/**
 * File renders a `ListItem` for a regular file, also showing it's size
 * @param {Object} props - React props
 * @param {String} props.name - Name of the file
 * @param {Number} props.size - Size in bytes
 * @param {string} props.path - Full path of the file in the repository
 */
const File = ({ name, size, path }) => {
  const { job, snapshot } = useParams();
  return (
    <ItemIconLink
      name={name}
      text={{ secondary: formatBytes(size) }}
      to={`${addr}/jobs/${job}/dump?path=${path}&snapshot=${snapshot}`}
      external
      icon={<InsertDriveFile />}
    />
  );
};

/**
 * ItemIconLink is a `ListItem` that also renders a primary text, `Link` and `Icon`
 * @param {Object} props - React props
 * @param {String} props.name - Primary text
 * @param {Object} props.text - Additional `props` for the primary text
 * @param {String} props.to - `href` for the `Link`
 * @param {Boolean} props.external - Whether to use `<a>` instead of `<Link>`
 * @param {JSX.Element} props.children - Optional component to also include
 */
const ItemIconLink = ({ name, text, to, external, icon, children }) => (
  <ListItem button component={external ? "a" : Link} to={to} href={to}>
    <ListItemIcon>{icon}</ListItemIcon>
    <ListItemText primary={name} {...text} />
    {children}
  </ListItem>
);

