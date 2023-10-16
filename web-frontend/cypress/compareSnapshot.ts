// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
const { VITE_CI_TEST } = import.meta.env || Cypress.env();

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function compareSnapshot(condition: any, name: string) {
  if (VITE_CI_TEST) {
    // originally is: cy.compareSnapshot('name-we-want');
    condition?.compareSnapshot(name);
  }
}
