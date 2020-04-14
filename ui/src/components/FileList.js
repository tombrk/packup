import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import axios from "axios";
import { addr } from "../api";

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
    <ul>
      {files.map((f) => (
        <li key={f.name}>
          {f.type === "dir" ? <Dir name={f.name} /> : f.name}
        </li>
      ))}
    </ul>
  );
};

/**
 * Directory renders a restic directory, including link
 * @param {Object} props - React props
 * @param {string} props.name - Name of the directory (no slashes)
 */
const Dir = ({ name }) => {
  const location = useLocation();
  return <Link to={`${location.pathname}/${name}`}>{name}</Link>;
};
