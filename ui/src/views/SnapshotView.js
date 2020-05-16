/** @jsx jsx */
import { jsx } from "@emotion/core";
import styled from "@emotion/styled";

import { useEffect, useState } from "react";
import { Link as RouterLink, useParams } from "react-router-dom";
import axios from "axios";
import { addr } from "../api";

import { Breadcrumbs, IconButton, Link, Paper } from "@material-ui/core";
import { ArrowBack } from "@material-ui/icons";
import { useTheme } from "@material-ui/core/styles";

import { Layout, AppTitle } from "./Layout";
import { FileList } from "../components/FileList";

import { SnapshotRoute } from "../UI";

/**
 * SnapshotView displays snapshot contents (list of files)
 */
export const SnapshotView = () => {
  const { job, snapshot, path } = useParams();
  const rPath = `/${useParams().path || ""}`;

  // path elements for navigation
  const elems = [snapshot, ...(path || "").split("/")];

  // load filelist from api
  const [files, setFiles] = useState([]);
  useEffect(() => {
    const fetch = async () => {
      setFiles([]);
      try {
        const result = await axios(`${addr}/jobs/${job}/files?path=${rPath}`);
        setFiles(result.data);
      } catch (error) {
        const message = error.response
          ? error.response.data
          : "Unable to list files. Please check your backend connection";

        console.error(message);
      }
    };
    fetch();
  }, [rPath]);

  return (
    <Layout
      preTitle={
        <div css={{ display: "flex", alignItems: "center" }}>
          <IconButton
            to="/"
            component={RouterLink}
            css={{ marginRight: useTheme().spacing(1) }}
            edge="start"
            color="inherit"
          >
            <ArrowBack />
          </IconButton>
          <AppTitle variant="h6">{job}</AppTitle>
        </div>
      }
      title={<PathNav elems={elems} job={job} />}
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
const PathNav = ({ elems, job }) => (
  <HeaderPaper>
    <Breadcrumbs>
      {elems.map((e, i) => (
        <Link
          component={RouterLink}
          key={i}
          to={`/${job}/${elems.slice(0, i + 1).join("/")}`}
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
