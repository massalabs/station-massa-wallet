import renderer from 'react-test-renderer';
import { cleanup } from '@testing-library/react';
import Body from './Body';

afterEach(cleanup);

describe('Body', () => {
  it('should render a body with the given text', () => {
    const message = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc consequat aliquam felis,
			ac vestibulum nisl
			placerat in. In hac habitasse platea dictumst. Pellentesque convallis erat eu enim tincidunt, a consectetur
			elit efficitur.`;
    const body = <Body>{message}</Body>;

    // snapshot test
    const componentSnapchat = renderer.create(body);
    const tree = componentSnapchat.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
