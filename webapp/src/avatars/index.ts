import avatar00 from "./avatar00.svg";
import avatar01 from "./avatar01.svg";
import avatar02 from "./avatar02.svg";
import avatar03 from "./avatar03.svg";
import avatar04 from "./avatar04.svg";
import avatar05 from "./avatar05.svg";
import avatar06 from "./avatar06.svg";
import avatar07 from "./avatar07.svg";
import avatar08 from "./avatar08.svg";

const avatars: Record<string, string> = {
  avatar00,
  avatar01,
  avatar02,
  avatar03,
  avatar04,
  avatar05,
  avatar06,
  avatar07,
  avatar08,
};

const avatar = (name: string): string | undefined => avatars[name];

export { avatar };
