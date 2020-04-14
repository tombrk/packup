import React from "react";
import { Link, useParams, useLocation } from "react-router-dom";
import { FileList } from "../components/FileList";

/**
 * SnapshotView displays snapshot contents (list of files)
 */
export const SnapshotView = () => {
  const path = `/${useParams().path || ""}`;

  const elems = useLocation().pathname.slice(1).split("/");

  return (
    <div>
      <PathNav elems={elems}></PathNav>
      <FileList path={path}></FileList>
    </div>
  );
};

/**
 * PathNav is the explorer-like navigation
 *
 * @param {Object} props - React props
 * @param {String[]} props.elems - current path splitted to individual folder names (`string.split('/')`)
 */
const PathNav = ({ elems }) => (
  <div>
    {elems.map((e, i) => (
      <div style={{ display: "inline" }} key={i}>
        <Link to={`/${elems.slice(0, i + 1).join("/")}`} key={i}>
          {e}
        </Link>
        <span>/</span>
      </div>
    ))}
  </div>
);
