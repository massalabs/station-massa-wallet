import { Meta, StoryObj } from '@storybook/react';
import AccountListItem from './AccountListItem';

const meta: Meta<typeof AccountListItem> = {
  title: 'AccountListItem',
  component: AccountListItem,
};

export default meta;

type Story = StoryObj<typeof AccountListItem>;

export const Primary: Story = {
  render: () => <AccountListItem name={'Account name 1'} amount={1000.11} />,
};
