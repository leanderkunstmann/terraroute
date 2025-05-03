import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useNavigate,
  useLocation,
} from 'react-router'
import * as React from 'react'
import AppBar from '@mui/material/AppBar'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import Box from '@mui/material/Box'
import Divider from '@mui/material/Divider'
import Drawer from '@mui/material/Drawer'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import MenuIcon from '@mui/icons-material/Menu'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import IconButton from '@mui/material/IconButton'

const drawerWidth = 240
const navItems = ['Globe', 'About', 'Contact']

export function Layout({ children }: { children: React.ReactNode }) {
  const [mobileOpen, setMobileOpen] = React.useState(false)
  const location = useLocation() // Get the current location
  const navigate = useNavigate()

  const theme = createTheme({
    palette: {
      mode: 'dark', //checked ? "dark" : "light",
    },
  })

  const handleDrawerToggle = () => {
    setMobileOpen((prevState) => !prevState)
  }

  const drawer = (
    <ThemeProvider theme={theme}>
      <Box onClick={handleDrawerToggle} sx={{ textAlign: 'center' }}>
        <Typography variant="h6" sx={{ my: 2 }}>
          Terraroute
        </Typography>
        <Divider />
        <List>
          {navItems.map((item) => (
            <ListItem key={item} disablePadding>
              <ListItemButton
                sx={{ textAlign: 'center' }}
                onClick={() => {
                  navigate('/' + item.toLowerCase())
                }}
              >
                <ListItemText primary={item} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Box>
    </ThemeProvider>
  )

  const container =
    typeof window !== 'undefined' ? () => window.document.body : undefined

  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <AppBar component="nav" color="inherit">
            <Toolbar variant="dense">
              <Typography
                variant="h6"
                component="div"
                sx={{ display: { xs: 'block', sm: 'flex' }, mr: 2 }}
                onClick={() => {
                  navigate('/')
                }}
                style={{ cursor: 'pointer' }}
              >
                Terraroute
              </Typography>
              <Box sx={{ flexGrow: 0, display: { xs: 'none', sm: 'flex' } }}>
                {navItems.map((item) => (
                  <Button
                    key={item}
                    onClick={() => {
                      navigate('/' + item.toLowerCase())
                    }}
                    size="large"
                    variant="text"
                    sx={{
                      color: 'text.primary',
                      backgroundColor:
                        location.pathname === '/' + item.toLowerCase()
                          ? 'action.selected'
                          : 'inherit',
                    }}
                  >
                    {item}
                  </Button>
                ))}
              </Box>
              <Box sx={{ flexGrow: 1 }} />
              <IconButton
                color="inherit"
                aria-label="open drawer"
                edge="end"
                onClick={handleDrawerToggle}
                sx={{ ml: 2, display: { sm: 'none' } }}
              >
                <MenuIcon />
              </IconButton>
            </Toolbar>
          </AppBar>

          <nav>
            <Drawer
              container={container}
              variant="temporary"
              anchor="right"
              open={mobileOpen}
              onClose={handleDrawerToggle}
              ModalProps={{
                keepMounted: true, // Better open performance on mobile.
              }}
              sx={{
                display: { xs: 'block', sm: 'none' },
                '& .MuiDrawer-paper': {
                  boxSizing: 'border-box',
                  width: drawerWidth,
                },
              }}
            >
              {drawer}
            </Drawer>
          </nav>
          {children}
        </ThemeProvider>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  )
}

export default function App() {
  return <Outlet />
}
