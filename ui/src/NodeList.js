/** @jsx jsx */
import { jsx } from "@emotion/core";

import { File, Dir, Placeholder } from "./Node";
import { List } from "@material-ui/core";

export const PlaceholderList = (props) => (
  <List>
    {props.items.map((item, index) => (
      <Placeholder key={index} />
    ))}
  </List>
);

const NodeList = (props) => (
  <List>
    {props.nodes
      .filter((item) => item.path !== props.dir)
      .map((item) => {
        const Node = item.type === "dir" ? Dir : File;
        return <Node key={item.path} {...item} />;
      })}
  </List>
);

export default NodeList;
