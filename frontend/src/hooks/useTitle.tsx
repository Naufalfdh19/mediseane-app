import { useEffect } from 'react'

import { TITLE_PREFIX } from '../const/const'

export default function useTitle(title: string) {
  useEffect(()=> {
    document.title = TITLE_PREFIX + title
  },[title])
}
