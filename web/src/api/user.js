import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/auth/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

export function getUserList(params) {
  return request({
    url: '/v1/user/list',
    method: 'get',
    params
  })
}

export function deleteUser(id) {
  return request({
    url: 'v1/user/delete',
    method: 'post',
    // headers: { 'Content-Type': 'application/json' },
    data: { 'id': id }
  })
}

export function addUser(data) {
  return request({
    url: 'v1/user/add',
    method: 'post',
    data
  })
}

export function updateUser(data) {
  return request({
    url: 'v1/user/update',
    method: 'post',
    data
  })
}
