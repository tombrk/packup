import React from "react";
import { Link as RouterLink, useParams, useLocation } from "react-router-dom";
import { FileList } from "../components/FileList";

import { Paper, Breadcrumbs, Link, styled } from "@material-ui/core";
import { Layout } from "./Layout";

/**
 * SnapshotView displays snapshot contents (list of files)
 */
export const SnapshotView = () => {
  const path = `/${useParams().path || ""}`;
  const elems = useLocation().pathname.slice(1).split("/");

  return (
    <Layout title={<PathNav elems={elems} />}>
      <FileList path={path}></FileList>
    </Layout>
  );
};

/**
 * PathNav is the explorer-like navigation
 *
 * @param {Object} props - React props
 * @param {String[]} props.elems - current path splitted to individual folder names (`string.split('/')`)
 */
const PathNav = ({ elems }) => (
  <HeaderPaper>
    <Breadcrumbs>
      {elems.map((e, i) => (
        <Link
          component={RouterLink}
          key={i}
          to={`/${elems.slice(0, i + 1).join("/")}`}
        >
          {e}
        </Link>
      ))}
    </Breadcrumbs>
  </HeaderPaper>
);

/**
 * HeaderPaper is a `Paper` specially styled,
 * so it renders nicely in the TitleBar
 * @param {Object} props - React props
 * @param {JSX.Element} props.children - Component to render
 */
const HeaderPaper = styled(Paper)({
  display: "flex",
  flexGrow: 1,
  alignItems: "center",
  minHeight: "3em",
  paddingLeft: "0.5em",
  paddingRight: "0.5em",
});
