import { useContext } from 'react'
import { AuthContext } from '../providers/app-provider'

export function useAuth() {
  return useContext(AuthContext)
}
