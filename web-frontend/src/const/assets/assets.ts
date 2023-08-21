export const MAS = 'MAS';

export interface IDeleteAssetsErrors {
  phrase?: string;
}

export const assetDeleteErrors = {
  success: 204,
  badRequest: 400,
  invalidAddress: 404,
  serverError: 500,
};

export const deleteConfirm = 'DELETE';

export const assetImportErrors = {
  success: 201,
  badRequest: 400,
  invalidAddress: 422,
  notFound: 404,
  serverError: 500,
};

export interface InputsErrors {
  address?: string;
}
