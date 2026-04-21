// @ts-nocheck;

import Icon from './Icon';

type Props = {
  name?: string;
  size?: number;
  color?: string;
  stroke?: string;
  onClick?: () => void
  cursor?: string;
  className?: string;
};

export const ArrowLeftGray = (props: Props) => <Icon {...props} name="arrow-left-gray" />;
export const Back5 = (props: Props) => <Icon {...props} name="back-5" />;
export const Close = (props: Props) => <Icon {...props} name="close" />;
export const Copy = (props: Props) => <Icon {...props} name="copy" />;
export const Load = (props: Props) => <Icon {...props} name="load" />;
export const Logout = (props: Props) => <Icon {...props} name="logout" />;
export const Loop = (props: Props) => <Icon {...props} name="loop" />;
export const Mic = (props: Props) => <Icon {...props} name="mic" />;
export const Pause = (props: Props) => <Icon {...props} name="pause" />;
export const Pfp = (props: Props) => <Icon {...props} name="pfp" />;
export const PlayFilled = (props: Props) => <Icon {...props} name="play-filled" />;
export const Play = (props: Props) => <Icon {...props} name="play" />;
export const Retry = (props: Props) => <Icon {...props} name="retry" />;
export const Skip5 = (props: Props) => <Icon {...props} name="skip-5" />;
export const SparkleSoft = (props: Props) => <Icon {...props} name="sparkle-soft" />;
export const Stop = (props: Props) => <Icon {...props} name="stop" />;
export const Translate = (props: Props) => <Icon {...props} name="translate" />;
export const Upload = (props: Props) => <Icon {...props} name="upload" />;
export const User = (props: Props) => <Icon {...props} name="user" />;