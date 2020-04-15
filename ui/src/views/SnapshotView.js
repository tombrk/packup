import { Breadcrumbs, IconButton, Link, Paper } from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { Link as RouterLink, useLocation, useParams } from "react-router-dom";

import { ArrowBack } from "@material-ui/icons";
import { FileList } from "../components/FileList";
import { Layout } from "./Layout";
import { addr } from "../api";
import axios from "axios";
/** @jsx jsx */
import { jsx } from "@emotion/core";
import styled from "@emotion/styled";
import { useTheme } from "@material-ui/core/styles";

/**
 * SnapshotView displays snapshot contents (list of files)
 */
export const SnapshotView = () => {
  const path = `/${useParams().path || ""}`;

  // path elements for navigation
  const elems = useLocation().pathname.slice(1).split("/");

  // files for FileList
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

  return (
    <Layout
      preTitle={
        <IconButton
          to="/"
          component={RouterLink}
          css={{ marginRight: useTheme().spacing(2) }}
          edge="start"
          color="inherit"
        >
          <ArrowBack />
        </IconButton>
      }
      title={<PathNav elems={elems} />}
      loading={!files.length}
    >
      <FileList files={files}></FileList>
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
