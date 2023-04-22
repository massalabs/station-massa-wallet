import Banner, { BannerProps } from './Banner';

export default {
  component: Banner,
  title: 'Banner',
};

const Template = (args: BannerProps) => <Banner {...args} />;

export const DefaultBanner = () => (
  <Template>Welcome to My Awesome React TypeScript App!</Template>
);
