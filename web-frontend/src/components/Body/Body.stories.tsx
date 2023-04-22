import Body, { BodyProps } from './Body';

export default {
  component: Body,
  title: 'Body',
};

const Template = (args: BodyProps) => <Body {...args} />;

export const DefaultBody = () => (
  <Template>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc consequat
    aliquam felis, ac vestibulum nisl placerat in. In hac habitasse platea
    dictumst. Pellentesque convallis erat eu enim tincidunt, a consectetur elit
    efficitur.
  </Template>
);
