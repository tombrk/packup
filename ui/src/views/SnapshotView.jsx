import { useEffect, useState } from "react";
import { Link as RouterLink, useParams } from "react-router-dom";
import axios from "axios";
import { addr } from "../api";

import { Breadcrumbs, IconButton, Link, Paper, styled } from "@material-ui/core";
import { ArrowBack } from "@material-ui/icons";
import { useTheme } from "@material-ui/core/styles";

import { useSnackbar } from "notistack";

import { Layout, AppTitle } from "./Layout";
import { FileList } from "../components/FileList";

import SnapPicker from "../SnapPicker";

/**
 * SnapshotView displays snapshot contents (list of files)
 */
export const SnapshotView = () => {
  const { job, snapshot, path } = useParams();
  const rPath = `/${useParams().path || ""}`;

  // path elements for navigation
  const elems = [snapshot, ...(path || "").split("/")];

  // component state
  const [files, setFiles] = useState([]);
  const [snapshots, setSnapshots] = useState([]);

  const { enqueueSnackbar } = useSnackbar();
  const theme = useTheme();

  // load filelist from api
  useEffect(() => {
    const fetch = async () => {
      setFiles([]);
      try {
        const result = await axios(
          `${addr}/jobs/${job}/files?snapshot=${snapshot}&path=${rPath}`
        );
        setFiles(result.data);
      } catch (error) {
        const message = error.response
          ? error.response.data
          : "Backend connection failed";

        enqueueSnackbar(`Loading files: ${message}`, { variant: "error" });
      }
    };
    fetch();
  }, [rPath, snapshot, job]);

  // load snapshots from api
  useEffect(() => {
    const fetch = async () => {
      setSnapshots([]);
      try {
        const result = await axios(`${addr}/jobs/${job}/snapshots`);
        setSnapshots(result.data);
      } catch (error) {
        const message = error.response
          ? error.response.data
          : "Backend connection failed";

        enqueueSnackbar(`Loading snapshots: ${message}`, { variant: "error" });
      }
    };
    fetch();
  }, [job]);

  return (
    <Layout
      preTitle={
        <div style={{ display: "flex", alignItems: "center" }}>
          <IconButton
            to="/"
            component={RouterLink}
            style={{ marginRight: theme.spacing(1) }}
            edge="start"
            color="inherit"
          >
            <ArrowBack />
          </IconButton>
          <AppTitle variant="h6">{job}</AppTitle>
        </div>
      }
      title={
        <div style={{ display: "flex", flexDirection: "row", flexGrow: 1 }}>
          <PathNav elems={elems} job={job} />
          <HeaderPaper>
            <SnapPicker
              job={job}
              path={path}
              snapshots={snapshots}
              current={snapshot}
            ></SnapPicker>
          </HeaderPaper>
        </div>
      }
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
  <HeaderPaper style={{ flexGrow: 1 }}>
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
  alignItems: "center",
  minHeight: "3em",
  paddingLeft: "0.5em",
  paddingRight: "0.5em",
  marginLeft: "0.5em",
});
