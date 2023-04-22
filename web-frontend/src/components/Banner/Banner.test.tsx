import renderer from 'react-test-renderer';
import { cleanup, render } from '@testing-library/react';
import Banner from './Banner';

afterEach(cleanup);

describe('Banner', () => {
  it('should render a banner with the given text', () => {
    const message = 'Welcome to My Awesome React TypeScript App!';
    const banner = <Banner>{message}</Banner>;

    // snapshot test
    const componentSnapchat = renderer.create(banner);
    const tree = componentSnapchat.toJSON();
    expect(tree).toMatchSnapshot();

    // DOM test
    const componentDom = render(banner);
    expect(componentDom.getByText(message)).toBeTruthy();
  });
});
