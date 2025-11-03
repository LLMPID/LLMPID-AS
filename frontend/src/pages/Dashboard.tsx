import { useState, useEffect, useRef } from "react";
import {
  Box,
  Button,
  VStack,
  Text,
  Flex,
  Collapsible,
  Textarea,
} from "@chakra-ui/react";

import {
  MenuContent,
  MenuItem,
  MenuRoot,
  MenuTrigger,
} from "@/components/ui/menu";
import { FaSortAmountDown, FaSortAmountUp } from "react-icons/fa";
import {
  QueryClient,
  QueryClientProvider,
  useMutation,
  useQuery,
} from "@tanstack/react-query";
import { fetchClassification, classifyText } from "@/api/classification";
import Header from "@/components/ui/Header";
const queryClient = new QueryClient();

function Dashboard() {
  const [text, setText] = useState("");
  const [classification, setClassification] = useState<string | null>(null);
  const [activeTicket, setActiveTicket] = useState<number | null>(null);

  const [sort, setSort] = useState("desc");
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  
  const { data, error, isLoading } = useQuery({
    queryKey: ["classification", page, limit, sort],
    queryFn: fetchClassification,
  });

  const addClassification = useMutation({
    mutationFn: classifyText,
    onSuccess: (newClassification) => {
      setClassification(newClassification.result);
      queryClient.invalidateQueries({ queryKey: ["classification"] });
      setText("");
      setActiveTicket(null);
    },
  });
  const prevValues = useRef({ limit, sort });

  useEffect(() => {
    setActiveTicket(null);

    if (
      prevValues.current.limit !== limit ||
      prevValues.current.sort !== sort
    ) {
      setPage(1);
    }

    prevValues.current = { limit, sort };
  }, [limit, sort, page]);

  return (
    <Box bg="white" minH="100vh">
      <Box bg="white" maxW="1100px" mx="auto" alignItems="center">
        <Header />
        <form
          onSubmit={(e) => {
            e.preventDefault();
            addClassification.mutate(text);
          }}
        >
          <VStack spaceY={4} alignItems="center" my={50}>
            <Textarea
              placeholder="Enter text"
              value={text}
              onChange={(e) => setText(e.target.value)}
              bg="white"
              borderColor="gray.300"
              _focus={{ borderColor: "blue.400" }}
              color="black"
              resize="none"
              height="150px"
              overflowY="auto"
            />

            <Button
              type="submit"
              colorScheme="blue"
              padding={5}
              fontSize="lg"
              fontWeight="bold"
              borderRadius="full"
              boxShadow="md"
              bgGradient="linear(to-r, blue.400, blue.500)"
              _hover={{
                bgGradient: "linear(to-r, blue.500, blue.600)",
                boxShadow: "lg",
              }}
              _active={{
                bgGradient: "linear(to-r, blue.600, blue.700)",
                boxShadow: "inner",
              }}
              loading={addClassification.isPending}
            >
              Classify
            </Button>
            {classification && (
              <Text
                fontSize="lg"
                color={
                  classification === "demo"
                    ? "gray.500"
                    : classification === "Normal"
                    ? "green.700"
                    : classification === "Injection"
                    ? "red.600"
                    : "black"
                }
                fontWeight="medium"
              >
                Classification: {classification}
              </Text>
            )}
            {addClassification.isError && (
              <Text color="red.500" fontSize="sm">
                {addClassification.error.message}
              </Text>
            )}
          </VStack>
        </form>
        <Flex
          fontSize="lg"
          mt={8}
          mb={4}
          fontWeight="medium"
          justify="space-between"
        >
          <Flex gap="3">
            <Text pt="1">Number of entries:</Text>
            <MenuRoot>
              <MenuTrigger asChild>
                <Button variant="outline" size="sm">
                  {limit}
                </Button>
              </MenuTrigger>
              <MenuContent>
                <MenuItem value="10" onClick={() => setLimit(10)}>
                  10
                </MenuItem>
                <MenuItem value="25" onClick={() => setLimit(25)}>
                  25
                </MenuItem>
                <MenuItem value="50" onClick={() => setLimit(50)}>
                  50
                </MenuItem>
              </MenuContent>
            </MenuRoot>
          </Flex>
          <Button onClick={() => setSort(sort === "desc" ? "asc" : "desc")}>
            {sort === "desc" ? <FaSortAmountDown /> : <FaSortAmountUp />}
          </Button>
        </Flex>

        <Flex direction="column" gap={3}>
          {isLoading && <Text>Loading...</Text>}
          {error && (
            <Text color="red.500" fontSize="sm">
              {error.message}
            </Text>
          )}
          {data?.slice().map((c: any, index: number) => (
            <Collapsible.Root
              key={c.id}
              open={activeTicket === index}
              unmountOnExit
            >
              <Collapsible.Trigger
                onClick={() =>
                  setActiveTicket((prev) => (prev === index ? null : index))
                }
                py="3"
                px="4"
                bg="white"
                borderWidth="2px"
                borderRadius="lg"
                borderColor={
                  c.result === "demo"
                    ? "gray.500"
                    : c.result === "Normal"
                    ? "green.700"
                    : c.result === "Injection"
                    ? "red.600"
                    : "black"
                }
                boxShadow="sm"
                _hover={{ bg: "blue.50", boxShadow: "md" }}
                width="100%"
                backgroundColor={
                  c.result === "demo"
                    ? "gray.50"
                    : c.result === "Normal"
                    ? "green.50"
                    : c.result === "Injection"
                    ? "red.50"
                    : "black"
                }
              >
                <Flex align="center" position="relative" width="100%">
                  <Text fontSize="md" color="blue.700" fontWeight="bold">
                    Id: {c.id}
                  </Text>
                  <Text
                    fontSize="md"
                    fontWeight="bold"
                    color={
                      c.result === "demo"
                        ? "gray.500"
                        : c.result === "Normal"
                        ? "green.700"
                        : c.result === "Injection"
                        ? "red.600"
                        : "black"
                    }
                    position="absolute"
                    left="50%"
                    transform="translateX(-50%)"
                  >
                    Result: {c.result}
                  </Text>
                  <Box ml="auto" textAlign="right">
                    <Text fontSize="sm" color="gray.500">
                      Source: {c.source_name}
                    </Text>
                    <Text fontSize="sm" color="gray.500">
                      Date:{" "}
                      {new Date(c.created_at).toLocaleTimeString("en-US", {
                        hour: "numeric",
                        minute: "numeric",
                        second: "numeric",
                        hour12: true,
                      })}
                    </Text>
                  </Box>
                </Flex>
              </Collapsible.Trigger>
              <Collapsible.Content>
                <Box mt={3} p={3} bg="gray.100" borderRadius="md">
                  <Text fontSize="sm" color="gray.700">
                    <strong>Request Text:</strong>
                    <Text whiteSpace="pre-wrap">{c.request_text}</Text>
                  </Text>
                </Box>
              </Collapsible.Content>
            </Collapsible.Root>
          ))}
        </Flex>
        <Flex my={5} justify="center" gap={4}>
          <Button
            onClick={() => setPage((prev) => Math.max(prev - 1, 1))}
            disabled={page === 1}
          >
            Prev
          </Button>
          <Text fontSize="md" color="gray.700" pt="6px">
            Page {page}
          </Text>
          <Button
            onClick={() => setPage((prev) => prev + 1)}
            disabled={!data || data.length < limit}
          >
            Next
          </Button>
        </Flex>
      </Box>
    </Box>
  );
}

export default function RootApp() {
  return (
    <QueryClientProvider client={queryClient}>
      <Dashboard />
    </QueryClientProvider>
  );
}
