import IconMap from "./assets/IconMap";

const EmptyIcon = () => <div />;

const Icon = ({
  name,
  size,
  color,
  cursor,
  onClick,
  stroke,
  className,
  ...rest
}: {
  name?: string;
  size?: number;
  color?: string;
  cursor?: string;
  stroke?: string;
  onClick?: () => void;
  className?: string;
}) => {
  // @ts-ignore
  const Icon = IconMap[name] || EmptyIcon;
  return (
    <Icon
      color={color}
      style={{ width: size, height: size, cursor: cursor }}
      onClick={onClick}
      className={className}
      {...rest}
    />
  );
};

export default Icon;
