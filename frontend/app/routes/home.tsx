import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import { useNavigate } from "react-router";

const cards = [
  {
    id: 1,
    title: "Shortest Path between two airports",
    description:
      "Generate the shortest path between two airports with opting out countries that are not allowed to fly over.",
    path: "/globe",
  },
  {
    id: 2,
    title: "Available flight connections",
    description:
      "Select or click on an airport to see available flight connections.",
    path: "/globe",
  },
  {
    id: 3,
    title: "Most efficient available flight paths",
    description:
      "Select two airports to get the most efficient available flight paths with stopovers. Filter by airlinegroups or countries.",
    path: "/globe",
  },
];

export default function Home() {
  const navigate = useNavigate();
  return (
    <div
      style={{
        width: "100%",
        height: "80vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Grid container spacing={2}>
        {cards.map((card, index) => (
          <Grid size={{ xs: 12, sm: 4 }} key={index}>
            <Card key={index} sx={{ height: "100%" }}>
              <CardActionArea
                onClick={() => {
                  navigate(card.path);
                }}
                sx={{
                  height: "100%",
                  backgroundColor: "action.selected",
                  "&:hover": {
                    backgroundColor: "action.selectedHover",
                  },
                }}
              >
                <CardContent sx={{ height: "100%" }}>
                  <Typography variant="h5" component="div">
                    {card.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {card.description}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </div>
  );
}
