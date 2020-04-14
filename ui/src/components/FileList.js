import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import axios from "axios";
import { addr } from "../api";

import { List, ListItem, ListItemIcon, ListItemText } from "@material-ui/core";
import { InsertDriveFile, Folder } from "@material-ui/icons";

/**
 *  FileList lists files from a path of a Restic snapshot
 * @param {Object} props - React props
 * @param {string} props.path - Snapshot subpath to display
 */
export const FileList = ({ path }) => {
  const [files, setFiles] = useState([]);

  // load filelist from api
  useEffect(() => {
    const fetch = async () => {
      setFiles([]);
      try {
        const result = await axios(`${addr}/files?path=${path}`);
        setFiles(result.data);
      } catch (error) {
        const message = error.response
          ? error.response.data
          : "Unable to list files. Please check your backend connection";

        console.error(message);
      }
    };
    fetch();
  }, [path]);

  // display as unordered list
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
const Dir = ({ name }) => (
  <ItemIconLink
    name={`${name}/`}
    to={`${useLocation().pathname}/${name}`}
    icon={<Folder />}
  />
);

/**
 * File renders a `ListItem` for a regular file, also showing it's size
 * @param {Object} props - React props
 * @param {String} props.name - Name of the file
 * @param {Number} props.size - Size in bytes
 * @param {string} props.path - Full path of the file in the repository
 */
const File = ({ name, size, path }) => (
  <ItemIconLink
    name={name}
    text={{ secondary: formatBytes(size) }}
    to={`${addr}/dump?path=${path}`}
    external
    icon={<InsertDriveFile />}
  />
);

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

/**
 * formatBytes converts a numeric byte count into a human readable unit (`KB`, `MB`, etc)
 * @param {Number} bytes - The bytes to format
 * @param {Number} decimals - Decimal count to use
 * @returns {String} Human readable string format
 */
function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}
