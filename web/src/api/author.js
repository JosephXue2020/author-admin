import request from '@/utils/request'

export function getAuthorList(params) {
  return request({
    url: '/v1/author/list',
    method: 'get',
    params
  })
}
